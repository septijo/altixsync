package main

import (
	//"bufio"
	//"bytes"
	"fmt"
	//"math/rand"
	//"net"
	"time"
)

//////////////////////////////////////////////////////////////////////////////
//                                                                          //
func send_data__device_setting__to__altix_device(slicePos int) bool {
	//                                                                          //
	//////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__device_setting__to__altix_device(", slicePos, ")")
	// 00, 07, yymdhns, Brightness, TotalCount_Gang (1-4), NotGood_Gang (1-4),
	//  Screen_Interval uint16, ScreenRow_JobShift[2][4] uint8, ScreenRow_Metric_Id[2][4] uint32
	data, deviceTimeStamp, jumlahId, err := MSSQL_get_device_setting(altixDevices[slicePos].altix_device_Id)
	if err != nil {
		return false
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(data)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt), len(data), "bytes sent")
	if writeCount != len(data) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	// ada di altixProtocol.go
	if !cekReply("send_data__device_setting__to__altix_device", 0x01000700, jumlahId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_device_setting_timeStamp = np_get_yymdhns_from_byte_slice(deviceTimeStamp)
	fmt.Println("altixDevices[slicePos].device_setting_timeStamp =", altixDevices[slicePos].dvc_device_setting_timeStamp)
	return true
} // end of send_data__device_setting__to__altix_device

//////////////////////////////////////////////////////////////////////////////
//                                                                          //
func send_data__master_metric__to__altix_device(slicePos int) bool {
	//                                                                          //
	//////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_metric__to__altix_device(", slicePos, ")")
	metricIds, mstMetrics, listErr, dataErr := MSSQL_get_mstMetrics(altixDevices[slicePos].tenant_Id,
		altixDevices[slicePos].dvc_master_metric_timeStamp)
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device dataErr:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstMetrics) == 0 {
		//fmt.Println("send_data__master_metric__to__altix_device ANEH .. jRec = 0 ??")
		//time.Sleep(5 * time.Second)
		//return true // harus return true agar ga ngaco
	}

	// 00, 08, jList, jData
	pjg := 60
	d2s := make([]byte, 4+4*len(metricIds)+pjg*len(mstMetrics))
	d2s[0] = 0
	d2s[1] = 8
	d2s[2] = byte(len(metricIds))
	d2s[3] = byte(len(mstMetrics))

	// bagian pertama adalah list semua metricId yang ada
	var jumlahListId uint32
	for i := 0; i < len(metricIds); i++ {
		jumlahListId += metricIds[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(metricIds[i], &d2s, i*4+4)
	}

	// bagian kedua adalah list semua data yang lebih baru dari tgl tersebut
	var jumlahDataId uint32
	var lastTimeStamp string
	startPos := 4 + 4*len(metricIds)
	for i := 0; i < len(mstMetrics); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += mstMetrics[i].metric_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstMetrics[i].metric_Id, &d2s, startPos+i*pjg+0)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstMetrics[i].job_limit_1_x_1K, &d2s, startPos+i*pjg+4)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstMetrics[i].job_limit_2_x_1K, &d2s, startPos+i*pjg+8)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstMetrics[i].shift_limit_1_x_1K, &d2s, startPos+i*pjg+12)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstMetrics[i].shift_limit_2_x_1K, &d2s, startPos+i*pjg+16)

		d2s[startPos+i*pjg+20] = byte(mstMetrics[i].Metric_Number)

		lastTimeStamp = mstMetrics[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+21)

		//for s := 0; s < 15; s++ {
		//	d2s[i*pjg+28+s] = mstMetrics[i].metric_Desc[s]
		//}
		//d2s[i*pjg+43] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		mdaX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(mstMetrics[i].display_As, 15)

		for s := 0; s < 15; s++ {
			d2s[startPos+i*pjg+28+s] = mdaX[s]
		}
		d2s[startPos+i*pjg+43] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		d2s[startPos+i*pjg+44] = mstMetrics[i].job_label_color
		d2s[startPos+i*pjg+45] = byte(mstMetrics[i].shift_label_color)

		d2s[startPos+i*pjg+46] = byte(mstMetrics[i].job_value_color_1)
		d2s[startPos+i*pjg+47] = byte(mstMetrics[i].job_value_color_2)
		d2s[startPos+i*pjg+48] = byte(mstMetrics[i].job_value_color_3)

		d2s[startPos+i*pjg+49] = byte(mstMetrics[i].shift_value_color_1)
		d2s[startPos+i*pjg+50] = byte(mstMetrics[i].shift_value_color_2)
		d2s[startPos+i*pjg+51] = byte(mstMetrics[i].shift_value_color_3)

		d2s[startPos+i*pjg+52] = mstMetrics[i].job_limit_type[0]
		d2s[startPos+i*pjg+53] = mstMetrics[i].shift_limit_type[0]

		// +54 dan +55 (2 byte) selanjutnya adalah untuk 32bit alignment
		d2s[startPos+i*pjg+54] = 0
		d2s[startPos+i*pjg+55] = 0

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4D4D5452, &d2s, startPos+i*pjg+56)
	}
	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_metric__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_metric__to__altix_device", 0x01000800, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].dvc_master_metric_timeStamp = lastTimeStamp
	return true
} // end of send_data__master_metric__to__altix_device
