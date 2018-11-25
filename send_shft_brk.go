package main

import (
	//"bufio"
	//"bytes"
	"fmt"
	//"math/rand"
	//"net"
	"time"
)

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
func send_data__master_shift__to__altix_device(slicePos int) bool {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_shift__to__altix_device(", slicePos, ")")
	// Down : 00, 0e, 1 byte jRow
	//    4 byte Shift_Id (uint32), 7 byte yymdhns, 1 byte Shift_Number,
	//    1 byte Color, 3 byte dummy_for_alignment
	lastDelete_dateTime, idS, mstShifts, listErr, dataErr := MSSQL_get_mstShifts(
		altixDevices[slicePos].tenant_Id, altixDevices[slicePos].dvc_master_shift_timeStamp[0])
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_Shift__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstShifts) == 0 {
		//fmt.Println("send_data__master_shift__to__altix_device ANEH .. jRec = 0 ??")
		//time.Sleep(5 * time.Second)
		//return true // harus return true agar ga ngaco
	}
	pjg := 20
	d2s := make([]byte, 4 + 7 + 4*len(idS)+pjg*len(mstShifts))
	d2s[0] = 0
	d2s[1] = 0x0E
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(mstShifts))

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
	for i := 0; i < len(mstShifts); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += mstShifts[i].Shift_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstShifts[i].Shift_Id, &d2s, startPos+i*pjg+0)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = mstShifts[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+4)

		d2s[startPos+i*pjg+11] = mstShifts[i].Shift_Number
		d2s[startPos+i*pjg+12] = mstShifts[i].Color

		//d2s[startPos + i*pjg+16] = mstShifts[i].Is_Deleted

		// lompat 3 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x53484654, &d2s, startPos+i*pjg+16)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_Shift__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_Shift__to__altix_device", 0x01000e00, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_master_shift_timeStamp[0] = lastTimeStamp
	altixDevices[slicePos].dvc_master_shift_timeStamp[1] = lastDelete_dateTime
	return true
} // end of send_data__master_shift__to__altix_device

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
func send_data__master_break__to__altix_device(slicePos int) bool {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_break__to__altix_device(", slicePos, ")")
	// Down : 00, 0F, 1 byte jRow
	//    4 byte Break_Id (uint32), 7 byte yymdhns, 1 byte Color,
	//    30 byte Break_Desc, 6 byte dummy_for_alignment
	lastDelete_dateTime, idS, mstBreaks, listErr, dataErr := MSSQL_get_mstBreaks(
		altixDevices[slicePos].tenant_Id, altixDevices[slicePos].dvc_master_break_timeStamp[0])
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_Break__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstBreaks) == 0 {
		//fmt.Println("send_data__master_break__to__altix_device ANEH .. jRec = 0 ??")
		//return true // harus return true agar ga ngaco
	}
	pjg := 48
	d2s := make([]byte, 4 + 7 + 4 * len(idS) + pjg * len(mstBreaks))
	d2s[0] = 0
	d2s[1] = 0x0F
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(mstBreaks))

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
	startPos := 4 + 7 + 4 * len(idS)
	for i := 0; i < len(mstBreaks); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += mstBreaks[i].Break_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstBreaks[i].Break_Id, &d2s, startPos + i*pjg+0)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = mstBreaks[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos + i*pjg+4)

		d2s[startPos + i*pjg+11] = mstBreaks[i].Color

		mbX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(mstBreaks[i].Break_Desc,30)
		for s := 0; s < 29; s++ {
			d2s[startPos + i*pjg+12+s] = mbX[s]
		}
		d2s[startPos + i*pjg+41] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		//d2s[startPos + i*pjg+45] = mstBreaks[i].Is_Deleted

		// lompat 2 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4D42524B, &d2s, startPos + i*pjg + 44)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_break__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_break__to__altix_device", 0x01000f00, jumlahListId + jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_master_break_timeStamp[0] = lastTimeStamp
	altixDevices[slicePos].dvc_master_break_timeStamp[1] = lastDelete_dateTime
	return true
} // end of send_data__master_break__to__altix_device
