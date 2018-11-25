package main

import (
	"bufio"
	//"bytes"
	"fmt"
	//"math/rand"
	//"net"
	"time"
)

///////////////////////////////////////////////////////////////////////////////////
//                                                                               //
func send_data__timeSchedule_Shift__to__altix_device(slicePos int) bool {
	//                                                                               //
	///////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__timeSchedule_Shift__to__altix_device(", slicePos, ")")

	lastDelete_dateTime, idS, timeSchedule_Shifts, listErr, dataErr := MSSQL_get_timeSchedule_Shifts(altixDevices[slicePos].machine_Id,
		altixDevices[slicePos].dvc_timeSchedule_Shift_timeStamp[0])
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__timeSchedule_Shift__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(timeSchedule_Shifts) == 0 {
		//fmt.Println("send_data__timeSchedule_shift__to__altix_device ANEH .. jRec = 0 ??")
		//return true // harus return true agar ga ngaco
	}
	pjg := 28
	d2s := make([]byte, 4 + 7 + 4 * len(idS) + pjg * len(timeSchedule_Shifts))
	d2s[0] = 0
	d2s[1] = 0x10
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(timeSchedule_Shifts))

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
	for i := 0; i < len(timeSchedule_Shifts); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += timeSchedule_Shifts[i].TimeSchedule_Shift_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(timeSchedule_Shifts[i].TimeSchedule_Shift_Id, &d2s, startPos + i*pjg+0)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(timeSchedule_Shifts[i].Shift_Id, &d2s, startPos + i*pjg+4)

		d2s[startPos + i*pjg+8] = byte(timeSchedule_Shifts[i].Shift_Minute_Duration % 256)
		d2s[startPos + i*pjg+9] = byte(timeSchedule_Shifts[i].Shift_Minute_Duration / 256)

		for ss := 0; ss < 5; ss++ { d2s[startPos + i*pjg + 10 + ss] = timeSchedule_Shifts[i].Shift_Start_HhMm[ss] }

		// ingat, dotNet, 0 = Minggu, 6 = Sabtu
        if timeSchedule_Shifts[i].Shift_Start_DayOfWeek > 6 {
			timeSchedule_Shifts[i].Shift_Start_DayOfWeek = 6
		}
		d2s[startPos + i*pjg+15] = byte(timeSchedule_Shifts[i].Shift_Start_DayOfWeek)
		
		d2s[startPos + i*pjg+16] = byte(timeSchedule_Shifts[i].CurrentOrNext)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = timeSchedule_Shifts[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos + i*pjg+17)

		//d2s[i*pjg+27] = timeSchedule_Shifts[i].Is_Deleted

		// lompat 0 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x54535348, &d2s, startPos + i*pjg + 24)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__timeSchedule_Shift__to__altix_device :", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__timeSchedule_Shift__to__altix_device", 0x01001000, jumlahListId + jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_timeSchedule_Shift_timeStamp[0] = lastTimeStamp
	altixDevices[slicePos].dvc_timeSchedule_Shift_timeStamp[1] = lastDelete_dateTime
	return true
} // end of send_data__timeSchedule_Shift__to__altix_device

////////////////////////////////////////////////////////////////////////////////////////
//                                                                                    //
func send_data__timeSchedule_ShiftBreak__to__altix_device(slicePos int) bool {
	//                                                                                    //
	////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__timeSchedule_ShiftBreak__to__altix_device(", slicePos, ")")
	// Down : 00, 11, 1 byte jRow
	//    4 byte Shift_Id (uint32), 7 byte yymdhns, 1 byte Shift_Number,
	//    1 byte Color, 3 byte dummy_for_alignment
	lastDelete_dateTime, idS, timeSchedule_ShiftBreaks, listErr, dataErr := MSSQL_get_timeSchedule_ShiftBreaks(altixDevices[slicePos].machine_Id,
		altixDevices[slicePos].dvc_timeSchedule_ShiftBreak_timeStamp[0])
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__timeSchedule_ShiftBreak__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(timeSchedule_ShiftBreaks) == 0 {
		//fmt.Println("send_data__timeSchedule_ShiftBreak__to__altix_device ANEH .. jRec = 0 ??")
		//return true // harus return true agar ga ngaco
	}

	// dibagi per sequence : 1 utk awal, 2 utk listId, 3 utk data

	// SEQUENCE 1 : awal
	d2s1 := make([]byte, 7 + 7)
	d2s1[0] = 0; d2s1[1] = 0x11; d2s1[2] = 1
	d2s1[3] = byte(len(idS) % 256);	                     d2s1[4] = byte(len(idS) / 256)
	d2s1[5] = byte(len(timeSchedule_ShiftBreaks) % 256); d2s1[6] = byte(len(timeSchedule_ShiftBreaks) / 256)

	fmt.Println(lastDelete_dateTime)
	np_write_datetime120_to_byte_slice7(lastDelete_dateTime, &d2s1, 7)

	fmt.Println(d2s1)
	writeCount, werr := altixDevices[slicePos].conn.Write(d2s1)
	if writeCount != len(d2s1) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__timeSchedule_ShiftBreak__SEQ_1_to__altix_device :", len(d2s1), "bytes sent")
	if !cekReply_TSSB(1, uint32(len(idS) + len(timeSchedule_ShiftBreaks)), slicePos) {
		return false
	}

	// SEQUENCE 2 : List ID
	maxRecList := 1390 / 4
	recListSent := 0
	for recListSent < len(idS) {
		jRec2Send := len(idS)
		if jRec2Send > maxRecList { jRec2Send = maxRecList}
		d2s2 := make([]byte,5 + jRec2Send*4)
		d2s2[0] = 0; d2s2[1] = 0x11; d2s2[2] = 2
		d2s2[3] = byte(jRec2Send % 256); d2s2[4] = byte(jRec2Send / 256)
		jumlahIdList := uint32(0)
		for rl := 0; rl < jRec2Send; rl++ {
			jumlahIdList += idS[recListSent + rl]
			nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[recListSent + rl], &d2s2, 5 + recListSent*4 + rl*4)
		}
		recListSent += jRec2Send
		writeCount, werr := altixDevices[slicePos].conn.Write(d2s2)
		if writeCount != len(d2s2) || werr != nil {
			// jangan close socket disini
			return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
		}
			fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__timeSchedule_ShiftBreak__SEQ_2_to__altix_device :", len(d2s2), "bytes sent")
		if !cekReply_TSSB(2, jumlahIdList, slicePos) {
			return false
		}
	}
	
	// SEQUENCE 3 : Data, asumsi pjg data TSSB = 32
	size := 32
	maxRecData := 1390 / size
	recDataSent := 0
	var lastTimeStamp string
	for recDataSent < len(timeSchedule_ShiftBreaks) {
		jRec2Send := len(timeSchedule_ShiftBreaks)
		if jRec2Send > maxRecData { jRec2Send = maxRecData}
		d2s3 := make([]byte,5 + jRec2Send*size)
		d2s3[0] = 0; d2s3[1] = 0x11; d2s3[2] = 3
		d2s3[3] = byte(jRec2Send % 256); d2s3[4] = byte(jRec2Send / 256)
		jumlahIdData := uint32(0)
		for rd := 0; rd < jRec2Send; rd++ {
			tssb := timeSchedule_ShiftBreaks[recDataSent + rd]
			jumlahIdData += tssb.TimeSchedule_ShiftBreak_Id
			startPos := 5 + (recDataSent + rd)*size
			nc_write_uint32_to_byte_slice__LSB_to_MSB(tssb.TimeSchedule_ShiftBreak_Id, &d2s3, startPos + 0)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(tssb.TimeSchedule_Shift_Id, &d2s3, startPos + 4)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(tssb.Break_Id, &d2s3, startPos + 8)
			d2s3[startPos + 12] = byte(tssb.Break_Minute_Duration % 256)
			d2s3[startPos + 13] = byte(tssb.Break_Minute_Duration / 256)
			for ss := 0; ss < 5; ss++ { d2s3[startPos + 14 + ss] = tssb.Break_Start_HhMm[ss] }
			d2s3[startPos + 19] = tssb.Break_Start_DayOfWeek
			d2s3[startPos + 20] = tssb.Current_Next_Occasional

			lastTimeStamp = tssb.str_time_stamp
			np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s3, startPos + 21)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(0x54534252, &d2s3, startPos + 28)
		}
		recDataSent += jRec2Send
		writeCount, werr := altixDevices[slicePos].conn.Write(d2s3)
		if writeCount != len(d2s3) || werr != nil {
			// jangan close socket disini
			return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
		}
		fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__timeSchedule_ShiftBreak__SEQ_3_to__altix_device :", len(d2s3), "bytes sent")
		if !cekReply_TSSB(3, jumlahIdData, slicePos) {
			return false
		}
	}

	altixDevices[slicePos].dvc_timeSchedule_ShiftBreak_timeStamp[0] = lastTimeStamp
	altixDevices[slicePos].dvc_timeSchedule_ShiftBreak_timeStamp[1] = lastDelete_dateTime
	return true
} // end of send_data__timeSchedule_ShiftBreak__to__altix_device

/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
func cekReply_TSSB(sequence uint8, cekJumlah uint32, slicePos int) bool {
	//                                                                             //
	/////////////////////////////////////////////////////////////////////////////////
	// reply nya harus jumlah dari metric_id yang dikirim
	altixDevices[slicePos].conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	readCount := 0
	message := make([]byte, 100)
	var err error
	readCount, err = bufio.NewReader(altixDevices[slicePos].conn).Read(message)
	if err != nil {
		fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+"  send TSSB wait4reply timeOut!!!!, readCount=", readCount)
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if readCount != 12 {
		fmt.Println("  send TSSB wait4reply Error !!!!, readCount=",
			readCount, "(must be 12 byte).")
		// jangan close socket disini
		return false
	}
	replyCode := uint32(message[0]) + uint32(message[1])<<8 +
		uint32(message[2])<<16 + uint32(message[3])<<24
	if replyCode != 0x01001100 {
		fmt.Println("  send TSSB wait4reply Error !!!!, replyCode=",
			replyCode,"(must be", 0x01001100, ").")
		// jangan close socket disini
		return false
	}

	replySequence := uint8(uint32(message[4]) + uint32(message[5])<<8 +
		uint32(message[6])<<16 + uint32(message[7])<<24)
	if replySequence != sequence {
		fmt.Println("  send TSSB wait4reply Error !!!!, replySequence=",
			replySequence,"(must be", sequence, ").")
		// jangan close socket disini
		return false
	}

	replyJumlahId := uint32(message[8]) + uint32(message[9])<<8 +
		uint32(message[10])<<16 + uint32(message[11])<<24
	if replyJumlahId != cekJumlah {
		fmt.Println("  send TSSB wait4reply Error !!!!, jumlahId=",
			replyJumlahId, "(must be", cekJumlah, ").")
		// jangan close socket disini
		return false
	} else {
		fmt.Println("  send TSSB wait4reply OK !!!!, sequence=",replySequence,"jumlahId=", replyJumlahId)
	}
	return true
} // end of cekReply_TSSB