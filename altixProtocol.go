package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	//"strconv"
	"time"
)

//////////////////////////////////////////////////////////////////////////
//                                                                      //
func handleRequest_altixDevice_sync(conn net.Conn, slicePos int) {
	//                                                                      //
	//////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" slice #", slicePos, conn.RemoteAddr(), "connected.")

	// begitu altix device konek, AltixSyn kirim 45 byte :
	// 00 + 4_byte_random_A + 4_byte_random_B + 32 byte SHA256(random_A + '5758')

	// AltixDevice akan cek apakah SHA256 thd random_A + '5758' tsb betul, jika betul
	// akan send 117 byte :
	// 00 + 32 byte SHA256(random_B + '8888') + 7 byte yymdhns + 12 byte mcu_id +
	//      16 byte storage_id + 7 byte setting's yymdhnss + 7 byte mstMetric yymdhns
	//      + 7 byte downReason's yymdhns + 7 byte standbyReason's yymdhns
	//      + 7 byte setupReason's yymdhns + 7 bytes mstShift's yymdhns +
	//      + 7 byte mstBreak's yymdhns

	// altixSync akan cek apakah sha tsb betul, jika betul akan send data2

	var data2send = make([]byte, 41)
	var rnd = rand.Int63()
	fmt.Printf("%16x", rnd)

	data2send[0] = 0
	for i := 0; i < 8; i++ {
		data2send[i+1] = byte((rnd << uint(i*8)) >> 56)
	}
	fmt.Println(data2send[1:9])

	s := npSha256_salted(data2send[1:5], "5758")
	for i := 0; i < 32; i++ {
		data2send[i+9] = s[i]
	}

	conn.Write(data2send)

	var message = make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var readCount, err = bufio.NewReader(conn).Read(message)
	if err != nil {
		fmt.Println("hail timeOut!!!!, readCount=", readCount, ". Closing socket...")
		delay_and_close_socket(500, slicePos)
		return
	}

	// 0x00 sha256 yymdhns cpuid storid last_yymdhns
	mustRcv := 208; //145
	if readCount != mustRcv {
		fmt.Println("hail Error !!!!, readCount=", readCount, "(must be", mustRcv, "byte). Closing socket...")
		delay_and_close_socket(500, slicePos)
		return
	}
	if message[0] != 0 {
		fmt.Println("hail Error !!!!, first byte=", message[0], "(must be 0x00 byte). Closing socket...")
		delay_and_close_socket(500, slicePos)
		return
	}

	j := npSha256_salted(data2send[5:9], "8888")
	// message[1:33] harus SHA256 (Random_B + '8888')
	if bytes.Compare(message[1:33], j) != 0 {
		fmt.Println("hail Error !!!!, byte[1:33]=", message[1:33],
			"(must be sha256(randomB+'8888')=", j, "). Closing socket...")
		delay_and_close_socket(500, slicePos)
		return
	}

	mcu_Id := fmt.Sprintf("%024X", message[40:52])
	stor_Id := fmt.Sprintf("%032X", message[52:68])
	altix_rtc := np_get_yymdhns_from_byte_slice(message[33:40])

	dvc_device_setting_timeStamp := np_get_yymdhns_from_byte_slice(message[68:75])
	dvc_master_metric_timeStamp := np_get_yymdhns_from_byte_slice(message[75:82])

	/*
	master_downReason_timeStamp := np_get_yymdhns_from_byte_slice(message[82:89])
	master_standbyReason_timeStamp := np_get_yymdhns_from_byte_slice(message[89:96])
	master_setupReason_timeStamp := np_get_yymdhns_from_byte_slice(message[96:103])

	master_shift_timeStamp := np_get_yymdhns_from_byte_slice(message[103:110])
	master_break_timeStamp := np_get_yymdhns_from_byte_slice(message[110:117])

	timeSchedule_Shift_timeStamp := np_get_yymdhns_from_byte_slice(message[117:124])
	timeSchedule_ShiftBreak_timeStamp := np_get_yymdhns_from_byte_slice(message[124:131])

	master_product_timeStamp := np_get_yymdhns_from_byte_slice(message[131:138])
	job_timeStamp := np_get_yymdhns_from_byte_slice(message[138:145])
	*/

	dvc_master_downReason_timeStamp := np_get_array_yymdhns_from_byte_slice(message[82:96])
	dvc_master_standbyReason_timeStamp := np_get_array_yymdhns_from_byte_slice(message[96:110])
	dvc_master_setupReason_timeStamp := np_get_array_yymdhns_from_byte_slice(message[110:124])

	dvc_master_shift_timeStamp := np_get_array_yymdhns_from_byte_slice(message[124:138])
	dvc_master_break_timeStamp := np_get_array_yymdhns_from_byte_slice(message[138:152])

	dvc_timeSchedule_Shift_timeStamp := np_get_array_yymdhns_from_byte_slice(message[152:166])
	dvc_timeSchedule_ShiftBreak_timeStamp := np_get_array_yymdhns_from_byte_slice(message[166:180])

	dvc_master_product_timeStamp := np_get_array_yymdhns_from_byte_slice(message[180:194])
	dvc_job_timeStamp := np_get_array_yymdhns_from_byte_slice(message[194:208])

	fmt.Println("hail ok !",
		"\n     Altix RTC . . . . . . . . . . . =", altix_rtc,
		"\n     Altix's Setting   . . . . . . . =", dvc_device_setting_timeStamp,
		"\n     Altix's mstMetric   . . . . . . =", dvc_master_metric_timeStamp,
		"\n     Altix's mstDownReason   . . . . =", dvc_master_downReason_timeStamp,
		"\n     Altix's mstStandbyReason  . . . =", dvc_master_standbyReason_timeStamp,
		"\n     Altix's mstSetupReason  . . . . =", dvc_master_setupReason_timeStamp,
		"\n     Altix's mstShift  . . . . . . . =", dvc_master_shift_timeStamp,
		"\n     Altix's mstBreak  . . . . . . . =", dvc_master_break_timeStamp,
		"\n     Altix's timeSchedule_Shift  . . =", dvc_timeSchedule_Shift_timeStamp,
		"\n     Altix's timeSchedule_ShiftBreak =", dvc_timeSchedule_ShiftBreak_timeStamp,
		"\n     Altix's mstProduct  . . . . . . =", dvc_master_product_timeStamp,
		"\n     Altix's job . . . . . . . . . . =", dvc_job_timeStamp)

	fmt.Println("     McuID=", mcu_Id, "StorId=", stor_Id)

	if stor_Id == "00000000000000000000000000000000" {
		fmt.Println("Storage Id empty, rejected !!")
		delay_and_close_socket(5800, slicePos)
		return
	}

	tenant_Id, altix_device_Id, machine_Id, timeStamps, err :=
		MSSQL_get__tenantId__deviceId__timeStampS__of_mcuId(mcu_Id, stor_Id)
	if err != nil {
		fmt.Println(err.Error())
		delay_and_close_socket(5800, slicePos)
		return
	}

	// tepat sebelum hearbeat, baru ini diisi
	// setelah diisi altix_device_id nya, baru altix dotNet bisa notif
	//   (karena dotmet mencari slice berdasarkan altix_device_id)
	altixDevices[slicePos].slice_lock.Lock()
	altixDevices[slicePos].tenant_Id = tenant_Id
	altixDevices[slicePos].machine_Id = machine_Id
	altixDevices[slicePos].altix_device_Id = altix_device_Id
	altixDevices[slicePos].mcu_Id = mcu_Id

	altixDevices[slicePos].dvc_device_setting_timeStamp = dvc_device_setting_timeStamp
	altixDevices[slicePos].dvc_master_metric_timeStamp = dvc_master_metric_timeStamp

	altixDevices[slicePos].dvc_master_downReason_timeStamp = dvc_master_downReason_timeStamp
	altixDevices[slicePos].dvc_master_standbyReason_timeStamp = dvc_master_standbyReason_timeStamp
	altixDevices[slicePos].dvc_master_setupReason_timeStamp = dvc_master_setupReason_timeStamp

	altixDevices[slicePos].dvc_master_shift_timeStamp = dvc_master_shift_timeStamp
	altixDevices[slicePos].dvc_master_break_timeStamp = dvc_master_break_timeStamp

	altixDevices[slicePos].dvc_timeSchedule_Shift_timeStamp = dvc_timeSchedule_Shift_timeStamp
	altixDevices[slicePos].dvc_timeSchedule_ShiftBreak_timeStamp = dvc_timeSchedule_ShiftBreak_timeStamp

	altixDevices[slicePos].dvc_master_product_timeStamp = dvc_master_product_timeStamp
	altixDevices[slicePos].dvc_job_timeStamp = dvc_job_timeStamp

	altixDevices[slicePos].SQL_device_setting_timeStamp = timeStamps[0][0]
	altixDevices[slicePos].SQL_master_metric_timeStamp = timeStamps[1][0]

	altixDevices[slicePos].SQL_master_downReason_timeStamp = timeStamps[2]
	altixDevices[slicePos].SQL_master_standbyReason_timeStamp = timeStamps[3]
	altixDevices[slicePos].SQL_master_setupReason_timeStamp = timeStamps[4]

	altixDevices[slicePos].SQL_master_shift_timeStamp = timeStamps[5]
	altixDevices[slicePos].SQL_master_break_timeStamp = timeStamps[6]

	altixDevices[slicePos].SQL_timeSchedule_Shift_timeStamp = timeStamps[7]
	altixDevices[slicePos].SQL_timeSchedule_ShiftBreak_timeStamp = timeStamps[8]

	altixDevices[slicePos].SQL_master_product_timeStamp = timeStamps[9]
	altixDevices[slicePos].SQL_job_timeStamp = timeStamps[10]

	altixDevices[slicePos].slice_lock.Unlock()

	socket_perlu_diclose := false

	for {
		// max. 30 detik sudah harus reply "HB" (HeartBeat), kecuali ada notif dari asp.net
		i := 0
		for ; i < 6; i++ {
			// cek yg ada di slice, karena itu bisa di ubah  dari handle_dotnet
			// WARNING : jangan close socket disalam situ, harus kasih tau dgn cara
			//   set flag "socket_perlu_diclose"
			if cek_timeStamp_apakah__MSSQL_lebih_baru_dari_altixDevice(slicePos) {
				i = 0 // karena sudah kirim sesuatu, balik ke 30 detik lagi
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				// ada kemungkinan socket perlu di close karena "sesuatu error"
				altixDevices[slicePos].slice_lock.Lock()
				socket_perlu_diclose = altixDevices[slicePos].socket_perlu_diclose
				altixDevices[slicePos].slice_lock.Unlock()
				//fmt.Println("abcde", socket_perlu_diclose)
				if socket_perlu_diclose {
					break
				}
			}
			time.Sleep(5 * time.Second)
		}

		if socket_perlu_diclose {
			fmt.Println(time.Now().Format(ymdhnsDateTimeFmt) + " socket perlu di close gan !!")
			//kok ga pernah kesini
			closeConn(slicePos)
			return
		}

		// jika dalam 30 detik tidak ada yg perlu di send, harus send HeartBeat
		altixDevices[slicePos].slice_lock.Lock()
		_, err := conn.Write([]byte{255, 255}) // _ = writeCount
		fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" slice", slicePos, "send   'Heartbeat Ask'   (0xFF 0xFF)")
		altixDevices[slicePos].slice_lock.Unlock()

		if err != nil {
			fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" slice #", slicePos,
				conn.RemoteAddr(), " Error writing")
			closeConn(slicePos)
			return
		}

		// send heartbeat (FF FF) ke Altix_Device, harus di reply dgn FF FE
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		readCount, err = bufio.NewReader(conn).Read(message)
		if err != nil {
			fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" slice", slicePos, "HeartBeat timeOut!!!!, readCount=", readCount, "err=", err.Error(), " Closing socket...")
			delay_and_close_socket(500, slicePos)
			return
		}

		if readCount != 2 {
			fmt.Println("HeartBeat Error !!!!, readCount=", readCount, "(must be 2 byte). Closing socket...")
			delay_and_close_socket(500, slicePos)
			return
		}
		if message[0] != 0xFF || message[1] != 0xFE {
			fmt.Println("HeartBeat Error !!!!, recv=", message[0:2], "(must be FF FE). Closing socket...")
			delay_and_close_socket(500, slicePos)
			return
		}
		fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" slice", slicePos, "  recv 'Heartbeat Reply' (0xFF 0xFE)")
	}
} // end of handleRequest_altixDevice_sync

////////////////////////////////////////////////////////////////////////////////////
//                                                                                //
func cek_timeStamp_apakah__MSSQL_lebih_baru_dari_altixDevice(slicePos int) (
	adaYgLebihBaru_dan_sudah_dinotif_ke_altixDevice bool) {
	//                                                                                //
	////////////////////////////////////////////////////////////////////////////////////
	//fmt.Println("func cek_timeStamp_apakah__MSSQL_lebih_baru_dari_altixDevice(", slicePos, ")")
	altixDevices[slicePos].slice_lock.Lock()

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_device_setting_timeStamp := altixDevices[slicePos].SQL_device_setting_timeStamp
	dvc_device_setting_timeStamp := altixDevices[slicePos].dvc_device_setting_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_master_metric_timeStamp := altixDevices[slicePos].SQL_master_metric_timeStamp
	dvc_master_metric_timeStamp := altixDevices[slicePos].dvc_master_metric_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_master_downReason_timeStamp := altixDevices[slicePos].SQL_master_downReason_timeStamp
	dvc_master_downReason_timeStamp := altixDevices[slicePos].dvc_master_downReason_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_master_standbyReason_timeStamp := altixDevices[slicePos].SQL_master_standbyReason_timeStamp
	dvc_master_standbyReason_timeStamp := altixDevices[slicePos].dvc_master_standbyReason_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_master_setupReason_timeStamp := altixDevices[slicePos].SQL_master_setupReason_timeStamp
	dvc_master_setupReason_timeStamp := altixDevices[slicePos].dvc_master_setupReason_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_master_shift_timeStamp := altixDevices[slicePos].SQL_master_shift_timeStamp
	dvc_master_shift_timeStamp := altixDevices[slicePos].dvc_master_shift_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_master_break_timeStamp := altixDevices[slicePos].SQL_master_break_timeStamp
	dvc_master_break_timeStamp := altixDevices[slicePos].dvc_master_break_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_timeSchedule_Shift_timeStamp := altixDevices[slicePos].SQL_timeSchedule_Shift_timeStamp
	dvc_timeSchedule_Shift_timeStamp := altixDevices[slicePos].dvc_timeSchedule_Shift_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_timeSchedule_ShiftBreak_timeStamp := altixDevices[slicePos].SQL_timeSchedule_ShiftBreak_timeStamp
	dvc_timeSchedule_ShiftBreak_timeStamp := altixDevices[slicePos].dvc_timeSchedule_ShiftBreak_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	//SQL_master_product_timeStamp := altixDevices[slicePos].SQL_master_product_timeStamp
	//dvc_master_product_timeStamp := altixDevices[slicePos].dvc_master_product_timeStamp

	// harus di copy ke memVar dulu karena ini masih di Lock, habis ini di unlock
	SQL_job_timeStamp := altixDevices[slicePos].SQL_job_timeStamp
	dvc_job_timeStamp := altixDevices[slicePos].dvc_job_timeStamp

	altixDevices[slicePos].slice_lock.Unlock()

	if (SQL_device_setting_timeStamp <= dvc_device_setting_timeStamp) &&

		(SQL_master_metric_timeStamp <= dvc_master_metric_timeStamp) &&

		(SQL_master_downReason_timeStamp[0] <= dvc_master_downReason_timeStamp[0]) &&
		(SQL_master_downReason_timeStamp[1] <= dvc_master_downReason_timeStamp[1]) &&

		(SQL_master_standbyReason_timeStamp[0] <= dvc_master_standbyReason_timeStamp[0]) &&
		(SQL_master_standbyReason_timeStamp[1] <= dvc_master_standbyReason_timeStamp[1]) &&

		(SQL_master_setupReason_timeStamp[0] <= dvc_master_setupReason_timeStamp[0]) &&
		(SQL_master_setupReason_timeStamp[1] <= dvc_master_setupReason_timeStamp[1]) &&

		(SQL_master_shift_timeStamp[0] <= dvc_master_shift_timeStamp[0]) &&
		(SQL_master_shift_timeStamp[1] <= dvc_master_shift_timeStamp[1]) &&

		(SQL_master_break_timeStamp[0] <= dvc_master_break_timeStamp[0]) &&
		(SQL_master_break_timeStamp[1] <= dvc_master_break_timeStamp[1]) &&

		(SQL_timeSchedule_Shift_timeStamp[0] <= dvc_timeSchedule_Shift_timeStamp[0]) &&
		(SQL_timeSchedule_Shift_timeStamp[1] <= dvc_timeSchedule_Shift_timeStamp[1]) &&

		(SQL_timeSchedule_ShiftBreak_timeStamp[0] <= dvc_timeSchedule_ShiftBreak_timeStamp[0]) &&
		(SQL_timeSchedule_ShiftBreak_timeStamp[1] <= dvc_timeSchedule_ShiftBreak_timeStamp[1]) &&

		//(SQL_master_product_timeStamp[0] <= dvc_master_product_timeStamp[0]) &&
		//(SQL_master_product_timeStamp[1] <= dvc_master_product_timeStamp[1]) &&

		(SQL_job_timeStamp[0] <= dvc_job_timeStamp[0]) &&
		(SQL_job_timeStamp[1] <= dvc_job_timeStamp[1]) {
		return false
	}

	// di send satu persatu agar tcpBuffer nya ga kelolodan

	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_device_setting_timeStamp > dvc_device_setting_timeStamp {
		/////////////////////
		//
		//    AWAS !!
		//    kalau reply dari Send Device Setting tdk sesuai, jangan lanjut send data lain
		//    karena bisa saja Storage nya punya machineId yg beda
		//
		////////////////////
		return send_data__device_setting__to__altix_device(slicePos) // snd_ds_mm.go
	}

	// master metric harus di send dulu, karena "screen setting" butuh data dari master metric
	if SQL_master_metric_timeStamp > dvc_master_metric_timeStamp {
		return send_data__master_metric__to__altix_device(slicePos) // snd_ds_mm.go
	}

	//fmt.Println("z")
	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_master_downReason_timeStamp[0] > dvc_master_downReason_timeStamp[0] ||
	   SQL_master_downReason_timeStamp[1] > dvc_master_downReason_timeStamp[1]  {
		return send_data__master_reason__to__altix_device("Down", slicePos) // snd_mr.go
	}

	//fmt.Println("y")
	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_master_standbyReason_timeStamp[0] > dvc_master_standbyReason_timeStamp[0] ||
	   SQL_master_standbyReason_timeStamp[1] > dvc_master_standbyReason_timeStamp[1] {
			return send_data__master_reason__to__altix_device("Standby", slicePos) // snd_mr.go
	}

	//fmt.Println("x")
	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_master_setupReason_timeStamp[0] > dvc_master_setupReason_timeStamp[0] ||
	   SQL_master_setupReason_timeStamp[1] > dvc_master_setupReason_timeStamp[1] {
		return send_data__master_reason__to__altix_device("Setup", slicePos) // snd_mr.go
	}

	//fmt.Println("w")
	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_master_shift_timeStamp[0] > dvc_master_shift_timeStamp[0] ||
	   SQL_master_shift_timeStamp[1] > dvc_master_shift_timeStamp[1] {
		return send_data__master_shift__to__altix_device(slicePos) // snd_ms_mb.go
	}

	//fmt.Println("v")
	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_master_break_timeStamp[0] > dvc_master_break_timeStamp[0] ||
	   SQL_master_break_timeStamp[1] > dvc_master_break_timeStamp[1] {
		return send_data__master_break__to__altix_device(slicePos)
	}

	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_timeSchedule_Shift_timeStamp[0] > dvc_timeSchedule_Shift_timeStamp[0] ||
	   SQL_timeSchedule_Shift_timeStamp[1] > dvc_timeSchedule_Shift_timeStamp[1] {
		return send_data__timeSchedule_Shift__to__altix_device(slicePos)
	}

	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_timeSchedule_ShiftBreak_timeStamp[0] > dvc_timeSchedule_ShiftBreak_timeStamp[0] ||
	   SQL_timeSchedule_ShiftBreak_timeStamp[1] > dvc_timeSchedule_ShiftBreak_timeStamp[1] {
		return send_data__timeSchedule_ShiftBreak__to__altix_device(slicePos)
	}

	/* 
	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_master_product_timeStamp[0] > dvc_master_product_timeStamp[0] ||
	   SQL_master_product_timeStamp[1] > dvc_master_product_timeStamp[1] {
		return send_data__master_product__to__altix_device(slicePos)
	} 
	*/

	// di send satu persatu agar tcpBuffer nya ga kelolodan
	if SQL_job_timeStamp[0] > dvc_job_timeStamp[0] || SQL_job_timeStamp[1] > dvc_job_timeStamp[1] {
		return send_data__job__to__altix_device(slicePos)
	}

	return true
} // end of cek_timeStamp_apakah__MSSQL_lebih_baru_dari_altixDevice

/////////////////////////////////////////////////////////////////////////////////
//                                                                             //
func cekReply(keterangan string, harusReplyCode uint32,
	cekJumlah uint32, slicePos int) bool {
	//                                                                             //
	/////////////////////////////////////////////////////////////////////////////////
	// reply nya harus jumlah dari metric_id yang dikirim
	altixDevices[slicePos].conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	readCount := 0
	message := make([]byte, 1024)
	var err error
	readCount, err = bufio.NewReader(altixDevices[slicePos].conn).Read(message)
	if err != nil {
		fmt.Println(keterangan+" wait4reply timeOut!!!!, readCount=", readCount)
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if readCount != 8 {
		fmt.Println(keterangan+" wait4reply Error !!!!, readCount=",
			readCount, "(must be 8 byte).")
		// jangan close socket disini
		return false
	}
	replyCode := uint32(message[0]) + uint32(message[1])<<8 +
		uint32(message[2])<<16 + uint32(message[3])<<24
	if replyCode != harusReplyCode {
		errDesc := ""
		for i := 0; i<len(altix_errorMessages); i++ {
			if replyCode == altix_errorMessages[i].error_Id {
				errDesc = altix_errorMessages[i].error_Desc
				break
			}
		}
		fmt.Println(keterangan+" wait4reply Error !!!!\n  replyCode=",
			fmt.Sprintf("0x%08X",replyCode), errDesc,"\n   must be  ", fmt.Sprintf("0x%08X",harusReplyCode))
		// jangan close socket disini
		return false
	}

	replyJumlahId := uint32(message[4]) + uint32(message[5])<<8 +
		uint32(message[6])<<16 + uint32(message[7])<<24
	if replyJumlahId != cekJumlah {
		fmt.Println(keterangan+" wait4reply Error !!!!, jumlahId=",
			replyJumlahId, "(must be", cekJumlah, ").")
		// jangan close socket disini
		return false
	} else {
		fmt.Println(keterangan+" wait4reply OK !!!!, jumlahId=", replyJumlahId)
	}
	return true
} // end of cekReply

////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                            //
func delay_and_set_flag_socket_perlu_di_close(delayMiliSecond int, slicePos int) bool {
	//                                                                                            //
	////////////////////////////////////////////////////////////////////////////////////////////////
	//fmt.Println("dasfspdc a")
	time.Sleep(time.Duration(delayMiliSecond) * time.Millisecond)
	//fmt.Println("dasfspdc b")
	altixDevices[slicePos].slice_lock.Lock()
	//fmt.Println("dasfspdc c")
	altixDevices[slicePos].socket_perlu_diclose = true
	//fmt.Println("dasfspdc d")
	altixDevices[slicePos].slice_lock.Unlock()
	//fmt.Println("dasfspdc e")
	return false
}

////////////////////////////////////////////////////////////////////////
//                                                                    //
func delay_and_close_socket(delayMiliSecond int, slicePos int) {
	//                                                                    //
	////////////////////////////////////////////////////////////////////////
	time.Sleep(time.Duration(delayMiliSecond) * time.Millisecond)
	closeConn(slicePos)
}

/*
/////////////////////////////////////////////////////////////////////////////
//                                                                         //
func send_data__master_break__to__altix_device(slicePos int) bool {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_break__to__altix_device(", slicePos, ")")
	// Down : 00, 0F, 1 byte jRow
	//    4 byte Break_Id (uint32), 7 byte yymdhns, 1 byte Color,
	//    30 byte Break_Desc, 6 byte dummy_for_alignment
	idS, mstBreaks, listErr, dataErr := MSSQL_get_mstBreaks(altixDevices[slicePos].tenant_Id,
		altixDevices[slicePos].master_break_timeStamp)
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_Break__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstBreaks) == 0 {
		fmt.Println("send_data__master_break__to__altix_device ANEH .. jRec = 0 ??")
		time.Sleep(5 * time.Second)
		return true // harus return true agar ga ngaco
	}
	pjg := 48
	d2s := make([]byte, 4+4*len(idS)+pjg*len(mstBreaks))
	d2s[0] = 0
	d2s[1] = 0x0F
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(mstBreaks))

	// bagian pertama adalah list semua Id yang ada
	var jumlahListId uint32
	for i := 0; i < len(idS); i++ {
		jumlahListId += idS[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[i], &d2s, i*4+4)
	}

	// bagian kedua adalah list semua data yang lebih baru dari tgl tersebut
	var jumlahDataId uint32
	var lastTimeStamp string
	startPos := 4 + 4*len(idS)
	for i := 0; i < len(mstBreaks); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += mstBreaks[i].Break_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstBreaks[i].Break_Id, &d2s, startPos+i*pjg+0)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = mstBreaks[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+4)

		d2s[startPos+i*pjg+11] = mstBreaks[i].Color

		mbX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(mstBreaks[i].Break_Desc, 30)
		for s := 0; s < 29; s++ {
			d2s[startPos+i*pjg+12+s] = mbX[s]
		}
		d2s[startPos+i*pjg+41] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		//d2s[startPos + i*pjg+45] = mstBreaks[i].Is_Deleted

		// lompat 2 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4D42524B, &d2s, startPos+i*pjg+44)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_break__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_break__to__altix_device", 0x01000f00, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].master_break_timeStamp = lastTimeStamp
	return true
} // end of send_data__master_break__to__altix_device

///////////////////////////////////////////////////////////////////////////////////
//                                                                               //
func send_data__timeSchedule_Shift__to__altix_device(slicePos int) bool {
	//                                                                               //
	///////////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__timeSchedule_Shift__to__altix_device(", slicePos, ")")

	idS, timeSchedule_Shifts, listErr, dataErr := MSSQL_get_timeSchedule_Shifts(altixDevices[slicePos].machine_Id,
		altixDevices[slicePos].timeSchedule_Shift_timeStamp)
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__timeSchedule_Shift__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(timeSchedule_Shifts) == 0 {
		fmt.Println("send_data__timeSchedule_shift__to__altix_device ANEH .. jRec = 0 ??")
		time.Sleep(5 * time.Second)
		return true // harus return true agar ga ngaco
	}
	pjg := 28
	d2s := make([]byte, 4+4*len(idS)+pjg*len(timeSchedule_Shifts))
	d2s[0] = 0
	d2s[1] = 0x10
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(timeSchedule_Shifts))

	// bagian pertama adalah list semua Id yang ada
	var jumlahListId uint32
	for i := 0; i < len(idS); i++ {
		jumlahListId += idS[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[i], &d2s, i*4+4)
	}

	// bagian kedua adalah list semua data yang lebih baru dari tgl tersebut
	var jumlahDataId uint32
	var lastTimeStamp string
	startPos := 4 + 4*len(idS)
	for i := 0; i < len(timeSchedule_Shifts); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		jumlahDataId += timeSchedule_Shifts[i].TimeSchedule_Shift_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(timeSchedule_Shifts[i].TimeSchedule_Shift_Id, &d2s, startPos+i*pjg+0)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(timeSchedule_Shifts[i].Shift_Id, &d2s, startPos+i*pjg+4)

		d2s[startPos+i*pjg+8] = byte(timeSchedule_Shifts[i].Shift_Minute_Duration % 256)
		d2s[startPos+i*pjg+9] = byte(timeSchedule_Shifts[i].Shift_Minute_Duration / 256)

		for ss := 0; ss < 5; ss++ {
			d2s[startPos+i*pjg+10+ss] = timeSchedule_Shifts[i].Shift_Start_HhMm[ss]
		}

		// ingat, dotNet, 0 = Minggu, 6 = Sabtu
		if timeSchedule_Shifts[i].Shift_Start_DayOfWeek > 6 {
			timeSchedule_Shifts[i].Shift_Start_DayOfWeek = 6
		}
		d2s[startPos+i*pjg+15] = byte(timeSchedule_Shifts[i].Shift_Start_DayOfWeek)

		d2s[startPos+i*pjg+16] = byte(timeSchedule_Shifts[i].CurrentOrNext)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = timeSchedule_Shifts[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+17)

		//d2s[i*pjg+27] = timeSchedule_Shifts[i].Is_Deleted

		// lompat 0 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x54535348, &d2s, startPos+i*pjg+24)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__timeSchedule_Shift__to__altix_device :", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__timeSchedule_Shift__to__altix_device", 0x01001000, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].timeSchedule_Shift_timeStamp = lastTimeStamp
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
	idS, timeSchedule_ShiftBreaks, listErr, dataErr := MSSQL_get_timeSchedule_ShiftBreaks(altixDevices[slicePos].machine_Id,
		altixDevices[slicePos].timeSchedule_ShiftBreak_timeStamp)
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__timeSchedule_ShiftBreak__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(timeSchedule_ShiftBreaks) == 0 {
		fmt.Println("send_data__timeSchedule_ShiftBreak__to__altix_device ANEH .. jRec = 0 ??")
		time.Sleep(5 * time.Second)
		return true // harus return true agar ga ngaco
	}

	// dibagi per sequence : 1 utk awal, 2 utk listId, 3 utk data

	// SEQUENCE 1 : awal
	d2s1 := make([]byte, 7)
	d2s1[0] = 0
	d2s1[1] = 0x11
	d2s1[2] = 1
	d2s1[3] = byte(len(idS) % 256)
	d2s1[4] = byte(len(idS) / 256)
	d2s1[5] = byte(len(timeSchedule_ShiftBreaks) % 256)
	d2s1[4] = byte(len(timeSchedule_ShiftBreaks) / 256)
	fmt.Println(d2s1)
	writeCount, werr := altixDevices[slicePos].conn.Write(d2s1)
	if writeCount != len(d2s1) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__timeSchedule_ShiftBreak__SEQ_1_to__altix_device :", len(d2s1), "bytes sent")
	if !cekReply_TSSB(1, uint32(len(idS)+len(timeSchedule_ShiftBreaks)), slicePos) {
		return false
	}

	// SEQUENCE 2 : List ID
	maxRecList := 1390 / 4
	recListSent := 0
	for recListSent < len(idS) {
		jRec2Send := len(idS)
		if jRec2Send > maxRecList {
			jRec2Send = maxRecList
		}
		d2s2 := make([]byte, 5+jRec2Send*4)
		d2s2[0] = 0
		d2s2[1] = 0x11
		d2s2[2] = 2
		d2s2[3] = byte(jRec2Send % 256)
		d2s2[4] = byte(jRec2Send / 256)
		jumlahIdList := uint32(0)
		for rl := 0; rl < jRec2Send; rl++ {
			jumlahIdList += idS[recListSent+rl]
			nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[recListSent+rl], &d2s2, 5+recListSent*4+rl*4)
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
		if jRec2Send > maxRecData {
			jRec2Send = maxRecData
		}
		d2s3 := make([]byte, 5+jRec2Send*size)
		d2s3[0] = 0
		d2s3[1] = 0x11
		d2s3[2] = 3
		d2s3[3] = byte(jRec2Send % 256)
		d2s3[4] = byte(jRec2Send / 256)
		jumlahIdData := uint32(0)
		for rd := 0; rd < jRec2Send; rd++ {
			tssb := timeSchedule_ShiftBreaks[recDataSent+rd]
			jumlahIdData += tssb.TimeSchedule_ShiftBreak_Id
			startPos := 5 + (recDataSent+rd)*size
			nc_write_uint32_to_byte_slice__LSB_to_MSB(tssb.TimeSchedule_ShiftBreak_Id, &d2s3, startPos+0)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(tssb.TimeSchedule_Shift_Id, &d2s3, startPos+4)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(tssb.Break_Id, &d2s3, startPos+8)
			d2s3[startPos+12] = byte(tssb.Break_Minute_Duration % 256)
			d2s3[startPos+13] = byte(tssb.Break_Minute_Duration / 256)
			for ss := 0; ss < 5; ss++ {
				d2s3[startPos+14+ss] = tssb.Break_Start_HhMm[ss]
			}
			d2s3[startPos+19] = tssb.Break_Start_DayOfWeek
			d2s3[startPos+20] = tssb.Current_Next_Occasional

			lastTimeStamp = tssb.str_time_stamp
			np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s3, startPos+21)
			nc_write_uint32_to_byte_slice__LSB_to_MSB(0x54534252, &d2s3, startPos+28)
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

	altixDevices[slicePos].timeSchedule_ShiftBreak_timeStamp = lastTimeStamp
	return true
} // end of send_data__timeSchedule_ShiftBreak__to__altix_device

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
func send_data__master_product__to__altix_device(slicePos int) bool {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt)+" func send_data__master_product__to__altix_device(", slicePos, ")")
	// Down : 00, 12, 1 byte jRow
	//    4 byte Shift_Id (uint32), 7 byte yymdhns, 1 byte Shift_Number,
	//    1 byte Color, 3 byte dummy_for_alignment
	idS, mstProducts, listErr, dataErr := MSSQL_get_mstProducts(altixDevices[slicePos].tenant_Id,
		altixDevices[slicePos].master_product_timeStamp)
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__master_Product__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(mstProducts) == 0 {
		fmt.Println("send_data__master_product__to__altix_device ANEH .. jRec = 0 ??")
		time.Sleep(5 * time.Second)
		return true // harus return true agar ga ngaco
	}
	pjg := 104
	d2s := make([]byte, 4+4*len(idS)+pjg*len(mstProducts))
	d2s[0] = 0
	d2s[1] = 0x12
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(mstProducts))

	// bagian pertama adalah list semua Id yang ada
	var jumlahListId uint32
	for i := 0; i < len(idS); i++ {
		jumlahListId += idS[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[i], &d2s, i*4+4)
	}

	// bagian kedua adalah list semua data yang lebih baru dari tgl tersebut
	var jumlahDataId uint32
	var lastTimeStamp string
	startPos := 4 + 4*len(idS)
	for i := 0; i < len(mstProducts); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima

		nc_write_uint64_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Pieces_x_1M, &d2s, startPos+i*pjg+0)
		nc_write_uint64_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Pieces_x_1M, &d2s, startPos+i*pjg+8)

		jumlahDataId += mstProducts[i].Product_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Product_Id, &d2s, startPos+i*pjg+16)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Hours_x_1M, &d2s, startPos+i*pjg+20)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Minutes_x_1M, &d2s, startPos+i*pjg+24)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Ideal_Cycle_Seconds_x_1M, &d2s, startPos+i*pjg+28)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Hours_x_1M, &d2s, startPos+i*pjg+32)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Minutes_x_1M, &d2s, startPos+i*pjg+36)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Takt_Seconds_x_1M, &d2s, startPos+i*pjg+40)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Scale_Total_Count_x_1M, &d2s, startPos+i*pjg+44)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Scale_Reject_Count_x_1M, &d2s, startPos+i*pjg+48)

		//// ga jadi .. uswutu harus lompat 4 byte dummy !!
		nc_write_uint32_to_byte_slice__LSB_to_MSB(mstProducts[i].Slow_Cycle__Treshold_Sequence_Display, &d2s, startPos+i*pjg+52)

		d2s[startPos+i*pjg+56] = byte(mstProducts[i].Pct_Slow_Cycle % 256)
		d2s[startPos+i*pjg+57] = byte(mstProducts[i].Pct_Slow_Cycle / 256)

		d2s[startPos+i*pjg+58] = byte(mstProducts[i].Pct_Small_Stop % 256)
		d2s[startPos+i*pjg+59] = byte(mstProducts[i].Pct_Small_Stop / 256)

		d2s[startPos+i*pjg+60] = byte(mstProducts[i].Pct_Full_Stop % 256)
		d2s[startPos+i*pjg+61] = byte(mstProducts[i].Pct_Full_Stop / 256)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = mstProducts[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+62)

		mpX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(mstProducts[i].Product_Name, 30)
		for s := 0; s < 29; s++ {
			d2s[startPos+i*pjg+69+s] = mpX[s]
		}
		d2s[startPos+i*pjg+98] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		//d2s[startPos + i*pjg+98] = mstProducts[i].Is_Deleted

		// lompat 1 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x50524F44, &d2s, startPos+i*pjg+100)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__master_Product__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__master_Product__to__altix_device", 0x01001200, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].master_product_timeStamp = lastTimeStamp
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
	idS, jobs, listErr, dataErr := MSSQL_get_jobs(altixDevices[slicePos].machine_Id,
		altixDevices[slicePos].job_timeStamp)
	if listErr != nil {
		fmt.Println("send_data__master_metric__to__altix_device listErr:\n" + listErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if dataErr != nil {
		fmt.Println("send_data__job__to__altix_device err:\n" + dataErr.Error())
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}
	if len(jobs) == 0 {
		fmt.Println("send_data__job__to__altix_device ANEH .. jRec = 0 ??")
		time.Sleep(5 * time.Second)
		return true // harus return true agar ga ngaco
	}
	pjg := 144
	d2s := make([]byte, 4+4*len(idS)+pjg*len(jobs))
	d2s[0] = 0
	d2s[1] = 0x13
	d2s[2] = byte(len(idS))
	d2s[3] = byte(len(jobs))

	// bagian pertama adalah list semua Id yang ada
	var jumlahListId uint32
	for i := 0; i < len(idS); i++ {
		jumlahListId += idS[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(idS[i], &d2s, i*4+4)
	}

	// bagian kedua adalah list semua data yang lebih baru dari tgl tersebut
	var jumlahDataId uint32
	var lastTimeStamp string
	startPos := 4 + 4*len(idS)
	for i := 0; i < len(jobs); i++ {
		// sudah diurutkan sesuai tgl, agar kalau terputus di tengah, altix_device sudah menyimpan tgl "terakhir" yg diterima
		nc_write_uint64_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Pieces_x_1M, &d2s, startPos+i*pjg+0)
		nc_write_uint64_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Pieces_x_1M, &d2s, startPos+i*pjg+8)

		jumlahDataId += jobs[i].Job_Id
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Job_Id, &d2s, startPos+i*pjg+16)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Product_Id, &d2s, startPos+i*pjg+20)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Hours_x_1M, &d2s, startPos+i*pjg+24)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Minutes_x_1M, &d2s, startPos+i*pjg+28)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Ideal_Cycle_Seconds_x_1M, &d2s, startPos+i*pjg+32)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Hours_x_1M, &d2s, startPos+i*pjg+36)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Minutes_x_1M, &d2s, startPos+i*pjg+40)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Takt_Seconds_x_1M, &d2s, startPos+i*pjg+44)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Scale_Total_Count_x_1M, &d2s, startPos+i*pjg+48)
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Scale_Reject_Count_x_1M, &d2s, startPos+i*pjg+52)

		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Goal_Qty, &d2s, startPos+i*pjg+56)

		//// ga jadi .. uswutu harus lompat 4 byte dummy !!
		nc_write_uint32_to_byte_slice__LSB_to_MSB(jobs[i].Slow_Cycle__Treshold_Sequence_Display, &d2s, startPos+i*pjg+60)

		d2s[startPos+i*pjg+64] = byte(jobs[i].Pct_Slow_Cycle % 256)
		d2s[startPos+i*pjg+65] = byte(jobs[i].Pct_Slow_Cycle / 256)

		d2s[startPos+i*pjg+66] = byte(jobs[i].Pct_Small_Stop % 256)
		d2s[startPos+i*pjg+67] = byte(jobs[i].Pct_Small_Stop / 256)

		d2s[startPos+i*pjg+68] = byte(jobs[i].Pct_Full_Stop % 256)
		d2s[startPos+i*pjg+69] = byte(jobs[i].Pct_Full_Stop / 256)

		// nomor "+ 3" s/d "+ 9" utk yymdhns, nomor "+ 10" utk padding alignment
		lastTimeStamp = jobs[i].str_time_stamp
		np_write_datetime120_to_byte_slice7(lastTimeStamp, &d2s, startPos+i*pjg+70)

		jdX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(jobs[i].Job_Desc, 30)
		for s := 0; s < 29; s++ {
			d2s[startPos+i*pjg+77+s] = jdX[s]
		}
		d2s[startPos+i*pjg+106] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		jrX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(jobs[i].Ref_No, 30)
		for s := 0; s < 29; s++ {
			d2s[startPos+i*pjg+107+s] = jrX[s]
		}
		d2s[startPos+i*pjg+136] = 0 // karena byte terakhir string di c/c++ harus chr(0)

		//d2s[startPos + i*pjg+136] = jobs[i].Is_Deleted

		// lompat 3 byte dummy

		nc_write_uint32_to_byte_slice__LSB_to_MSB(0x4A4F4253, &d2s, startPos+i*pjg+140)
	}

	writeCount, werr := altixDevices[slicePos].conn.Write(d2s)
	fmt.Println(time.Now().Format(ymdhnsDateTimeFmt)+" send_data__job__to__altix_device:", len(d2s), "bytes sent")
	if writeCount != len(d2s) || werr != nil {
		// jangan close socket disini
		return delay_and_set_flag_socket_perlu_di_close(500, slicePos)
	}

	if !cekReply("send_data__job__to__altix_device", 0x01001300, jumlahListId+jumlahDataId, slicePos) {
		return false
	}
	altixDevices[slicePos].job_timeStamp = lastTimeStamp
	return true
} // end of send_data__job__to__altix_device

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
			replyCode, "(must be", 0x01001100, ").")
		// jangan close socket disini
		return false
	}

	replySequence := uint8(uint32(message[4]) + uint32(message[5])<<8 +
		uint32(message[6])<<16 + uint32(message[7])<<24)
	if replySequence != sequence {
		fmt.Println("  send TSSB wait4reply Error !!!!, replySequence=",
			replySequence, "(must be", sequence, ").")
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
		fmt.Println("  send TSSB wait4reply OK !!!!, sequence=", replySequence, "jumlahId=", replyJumlahId)
	}
	return true
} // end of cekReply
*/