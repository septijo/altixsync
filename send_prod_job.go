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
func send_data__master_product__to__altix_device(slicePos int) bool {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_product__to__altix_device(", slicePos, ")")
	// Down : 00, 12, 1 byte jRow
	//    4 byte Shift_Id (uint32), 7 byte yymdhns, 1 byte Shift_Number,
	//    1 byte Color, 3 byte dummy_for_alignment
	lastDelete_dateTime, idS, mstProducts, listErr, dataErr := MSSQL_get_mstProducts(altixDevices[slicePos].tenant_Id,
		altixDevices[slicePos].dvc_master_product_timeStamp[0])
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_Product__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstProducts) == 0 {
		//fmt.Println("send_data__master_product__to__altix_device ANEH .. jRec = 0 ??")
		//return true // harus return true agar ga ngaco
	}
	pjg := 104
	d2s := make([]byte, 4 + 7 + 4 * len(idS) + pjg * len(mstProducts))
	d2s[0] = 0
	d2s[1] = 0x12
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(mstProducts))

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
	for i := 0; i < len(mstProducts); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		nc_write_uint64_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Pieces_x_1M, &d2s, startPos + i*pjg+0)
		nc_write_uint64_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Pieces_x_1M, &d2s, startPos + i*pjg+8)

		jumlahDataId += mstProducts[i].Product_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Product_Id, &d2s, startPos + i*pjg+16)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Hours_x_1M, &d2s, startPos + i*pjg+20)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Minutes_x_1M, &d2s, startPos + i*pjg+24)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Seconds_x_1M, &d2s, startPos + i*pjg+28)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Hours_x_1M, &d2s, startPos + i*pjg+32)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Minutes_x_1M, &d2s, startPos + i*pjg+36)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Seconds_x_1M, &d2s, startPos + i*pjg+40)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Scale_Total_Count_x_1K, &d2s, startPos + i*pjg+44)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Scale_Reject_Count_x_1K, &d2s,startPos +  i*pjg+48)

		//// ga jadi .. uswutu harus lompat 4 byte dummy !!
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Slow_Cycle__Treshold_Sequence_Display, &d2s,startPos +  i*pjg+52)

		d2s[startPos + i*pjg+56] = byte(mstProducts[i].Pct_Slow_Cycle % 256)
		d2s[startPos + i*pjg+57] = byte(mstProducts[i].Pct_Slow_Cycle / 256)

		d2s[startPos + i*pjg+58] = byte(mstProducts[i].Pct_Small_Stop % 256)
		d2s[startPos + i*pjg+59] = byte(mstProducts[i].Pct_Small_Stop / 256)

		d2s[startPos + i*pjg+60] = byte(mstProducts[i].Pct_Full_Stop % 256)
		d2s[startPos + i*pjg+61] = byte(mstProducts[i].Pct_Full_Stop / 256)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = mstProducts[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos + i*pjg+62)

		mpX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(mstProducts[i].Product_Name,30)
		for s := 0; s < 29; s++ {
			d2s[startPos + i*pjg+69+s] = mpX[s]
		}
		d2s[startPos + i*pjg+98] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		//d2s[startPos + i*pjg+98] = mstProducts[i].Is_Deleted

		// lompat 1 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x50524F44, &d2s, startPos + i*pjg + 100)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_Product__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_Product__to__altix_device", 0x01001200, jumlahListId + jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_master_product_timeStamp[0] = lastTimeStamp
	altixDevices[slicePos].dvc_master_product_timeStamp[1] = lastDelete_dateTime
	return true
} // end of send_data__master_product__to__altix_device

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
func send_data__job__to__altix_device(slicePos int) bool {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__job__to__altix_device(", slicePos, ")")
	// Down : 00, 13, 1 byte jRow
	//    4 byte Shift_Id (uint32), 7 byte yymdhns, 1 byte Shift_Number,
	//    1 byte Color, 3 byte dummy_for_alignment
	lastDelete_dateTime, idS, jobs, listErr, dataErr := MSSQL_get_jobs(altixDevices[slicePos].machine_Id,
		altixDevices[slicePos].dvc_job_timeStamp[0])
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__job__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(jobs) == 0 {
		//fmt.Println("send_data__job__to__altix_device ANEH .. jRec = 0 ??")
		//return true // harus return true agar ga ngaco
	}
	pjg := 160
	d2s := make([]byte, 4 + 7 + 4 * len(idS) + pjg * len(jobs))
	d2s[0] = 0
	d2s[1] = 0x13
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(jobs))

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
	lastTimeStamp := "1900-01-01 00:00:00" // direset dulu, siapa tau jRec nya kosong (sudah di delete semua)
	startPos := 4 + 7 + 4 * len(idS)
	for i := 0; i < len(jobs); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima
		nc_write_uint64_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Pieces_x_1M, &d2s, startPos + i*pjg+0)
		nc_write_uint64_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Pieces_x_1M, &d2s, startPos + i*pjg+8)

		jumlahDataId += jobs[i].Job_Id
		ndx := 16;
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Job_Id, &d2s, startPos + i*pjg+ndx); ndx += 4

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Product_Id, &d2s, startPos + i*pjg+ndx); ndx += 4

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Hours_x_1M, &d2s, startPos + i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Minutes_x_1M, &d2s, startPos + i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Seconds_x_1M, &d2s, startPos + i*pjg+ndx); ndx += 4

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Hours_x_1M, &d2s, startPos + i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Minutes_x_1M, &d2s, startPos + i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Seconds_x_1M, &d2s, startPos + i*pjg+ndx); ndx += 4

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Scale_Total_Count_x_1K, &d2s, startPos + i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Scale_Reject_Count_x_1K, &d2s, startPos + i*pjg+ndx); ndx += 4

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Goal_Qty, &d2s, startPos + i*pjg+ndx); ndx += 4

		//// ga jadi .. uswutu harus lompat 4 byte dummy !!
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Slow_Cycle__Treshold_Sequence_Display, &d2s,startPos +  i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Down_Rate_per_Hour, &d2s,startPos +  i*pjg+ndx); ndx += 4
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Run_Rate_per_Hour, &d2s,startPos +  i*pjg+ndx); ndx += 4

		//d2s[startPos + i*pjg+64] = byte(jobs[i].Pct_Slow_Cycle % 256); 
		//d2s[startPos + i*pjg+65] = byte(jobs[i].Pct_Slow_Cycle / 256)
		nc_write_uint16_to_byte_slice__LSB_to_MSB(jobs[i].Pct_Slow_Cycle, &d2s,startPos +  i*pjg+ndx); ndx += 2

		//d2s[startPos + i*pjg+66] = byte(jobs[i].Pct_Small_Stop % 256)
		//d2s[startPos + i*pjg+67] = byte(jobs[i].Pct_Small_Stop / 256)
		nc_write_uint16_to_byte_slice__LSB_to_MSB(jobs[i].Pct_Small_Stop, &d2s,startPos +  i*pjg+ndx); ndx += 2

		//d2s[startPos + i*pjg+68] = byte(jobs[i].Pct_Full_Stop % 256)
		//d2s[startPos + i*pjg+69] = byte(jobs[i].Pct_Full_Stop / 256)
		nc_write_uint16_to_byte_slice__LSB_to_MSB(jobs[i].Pct_Full_Stop, &d2s,startPos +  i*pjg+ndx); ndx += 2

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = jobs[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos + i*pjg+ndx); ndx += 7

		jdX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(jobs[i].Job_Desc,30)
		for s := 0; s < 29; s++ {
			d2s[startPos + i*pjg+ndx] = jdX[s]; ndx++
		}
		d2s[startPos + i*pjg+ndx] = 0; ndx++; // karena byte terakhir string di c/c++ harus chr(0)

		fmt.Println("before filling job_Ref ndx=",ndx)
		jrX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(jobs[i].Ref_No,30)
		for s := 0; s < 29; s++ {
			d2s[startPos + i*pjg+ndx] = jrX[s]; ndx++
		}
		d2s[startPos + i*pjg+ndx] = 0; ndx++; // karena byte terakhir string di c/c++ harus chr(0)

		//d2s[startPos + i*pjg+136] = jobs[i].Is_Deleted
		d2s[startPos + i*pjg+ndx] = jobs[i].Switch_from_Run_to_Down_based_on; ndx++
		d2s[startPos + i*pjg+ndx] = jobs[i].Switch_from_Down_to_Run_based_on; ndx++
		d2s[startPos + i*pjg+ndx] = jobs[i].Switch_from_Down_to_Run_Count; ndx++
		//d2s[startPos + i*pjg+ndx] = jobs[i].Switch_from_Down_to_Run_counter_of; ndx++
		d2s[startPos + i*pjg+ndx] = jobs[i].Rate_per_Hour_counter_of; ndx++

		d2s[startPos + i*pjg+ndx] = jobs[i].Remote_Code; ndx++

		// lompat 6 byte dummy 
		ndx += 6

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4A4F4253, &d2s, startPos + i*pjg + ndx)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__job__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__job__to__altix_device", 0x01001300, jumlahListId + jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_job_timeStamp[0] = lastTimeStamp
	altixDevices[slicePos].dvc_job_timeStamp[1] = lastDelete_dateTime
	return true
} // end of send_data__job__to__altix_device