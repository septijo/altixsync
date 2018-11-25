package main

import (
	//"bufio"
	//"bytes"
	"fmt"
	//"math/rand"
	//"net"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                              //
func send_data__master_reason__to__altix_device(Reason_Type string, slicePos int) bool {
	//                                                                                              //
	//////////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_"+
		Reason_Type+"_reason__to__altix_device(", slicePos, ")")
	var ts string
	if Reason_Type == "Down" {
		ts = altixDevices[slicePos].dvc_master_downReason_timeStamp[0]
	}
	if Reason_Type == "Standby" {
		ts = altixDevices[slicePos].dvc_master_standbyReason_timeStamp[0]
	}
	if Reason_Type == "Setup" {
		ts = altixDevices[slicePos].dvc_master_setupReason_timeStamp[0]
	}

	lastDelete_dateTime, idS, mstReasons, listErr, dataErr := 
		MSSQL_get_mstReasons(Reason_Type, altixDevices[slicePos].tenant_Id, ts)
		
	if listErr != nil {
		fmt.Println("send_data__master_" + Reason_Type + "_reason__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_" + Reason_Type + "_reason__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstReasons) == 0 {
		//fmt.Println("send_data__master_" + Reason_Type + "_reason__to__altix_device ANEH .. jRec = 0 ??")
		//time.Sleep(5 * time.Second)
		//return true // harus return true agar ga ngaco
	}

	// 00, xx, jList, jData
	var pjg int
	if Reason_Type == "Down" {
		pjg = 48
	}
	if Reason_Type == "Standby" {
		pjg = 52
	}
	if Reason_Type == "Setup" {
		pjg = 56
	}
	var harusReplyCode uint32
	d2s := make([]byte, 4+7+4*len(idS)+pjg*len(mstReasons))
	d2s[0] = 0
	if Reason_Type == "Down" {
		d2s[1] = 0x0b
		harusReplyCode = 0x01000b00
	}
	if Reason_Type == "Standby" {
		d2s[1] = 0x0c
		harusReplyCode = 0x01000c00
	}
	if Reason_Type == "Setup" {
		d2s[1] = 0x0d
		harusReplyCode = 0x01000d00
	}
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(mstReasons))

	fmt.Println(lastDelete_dateTime)
	np_write_datetime120_to_byte_slice7(lastDelete_dateTime, &d2s, 4)

	// bagian pertama adalah list semua Id yang ada
	var jumlahListId uint32
	for i := 0; i < len(idS); i++ {
		jumlahListId += idS[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[i], &d2s, i*4 + 4 + 7)
	}

	// bagian kedua adalah list semua data yang lebih baru dari tgl tersebut
	var jumlahDataId uint32
	var lastTimeStamp string
	startPos := 4 + 7 + 4*len(idS)
	for i := 0; i < len(mstReasons); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += mstReasons[i].Reason_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstReasons[i].Reason_Id, &d2s, startPos+i*pjg+0)

		lastTimeStamp = mstReasons[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+4)

		d2s[startPos+i*pjg+11] = mstReasons[i].Remote_Code

		mrdX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(mstReasons[i].Reason_Desc, 30)
		for s := 0; s < 29; s++ {
			d2s[startPos+i*pjg+12+s] = mrdX[s]
		}
		d2s[startPos+i*pjg+41] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		if Reason_Type == "Down" {
			nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4D444E52, &d2s, startPos+i*pjg+44)
		}
		if Reason_Type == "Standby" || Reason_Type == "Setup" {
			d2s[startPos+i*pjg+42] = byte(mstReasons[i].Green_Duration % 256)
			d2s[startPos+i*pjg+43] = byte(mstReasons[i].Green_Duration / 256)
			d2s[startPos+i*pjg+44] = byte(mstReasons[i].Yellow_Duration % 256)
			d2s[startPos+i*pjg+45] = byte(mstReasons[i].Yellow_Duration / 256)
			d2s[startPos+i*pjg+46] = byte(mstReasons[i].Reason_Duration % 256)
			d2s[startPos+i*pjg+47] = byte(mstReasons[i].Reason_Duration / 256)
		}
		if Reason_Type == "Standby" {
			nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4D534252, &d2s, startPos+i*pjg+48)
		}
		if Reason_Type == "Setup" {
			d2s[startPos+i*pjg+48] = byte(mstReasons[i].GoodCount_to_End_Setup % 256)
			d2s[startPos+i*pjg+49] = byte(mstReasons[i].GoodCount_to_End_Setup / 256)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4D535552, &d2s, startPos+i*pjg+52)
		}
		//d2s[startPos + i*pjg+53] = byte(mstReasons[i].Is_Deleted)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_"+Reason_Type+
		"Reason__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_"+Reason_Type+"Reason__to__altix_device",
		harusReplyCode, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	if Reason_Type == "Down" {
		altixDevices[slicePos].dvc_master_downReason_timeStamp[0] = lastTimeStamp
		altixDevices[slicePos].dvc_master_downReason_timeStamp[1] = lastDelete_dateTime
	}
	if Reason_Type == "Standby" {
		altixDevices[slicePos].dvc_master_standbyReason_timeStamp[0] = lastTimeStamp
		altixDevices[slicePos].dvc_master_standbyReason_timeStamp[1] = lastDelete_dateTime
	}
	if Reason_Type == "Setup" {
		altixDevices[slicePos].dvc_master_setupReason_timeStamp[0] = lastTimeStamp
		altixDevices[slicePos].dvc_master_setupReason_timeStamp[1] = lastDelete_dateTime
	}
	return true
} // end of send_data__master_reason__to__altix_device
