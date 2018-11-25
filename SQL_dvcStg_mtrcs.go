package main

import (
	//"time"

	//"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
	//"strings"
)

type altix_mstMetric struct {
	metric_Id, job_limit_1_x_1K, job_limit_2_x_1K, shift_limit_1_x_1K,
	shift_limit_2_x_1K uint32
	metric_Desc, display_As, str_time_stamp string
	job_label_color, shift_label_color, job_value_color_1,
	job_value_color_2, job_value_color_3, shift_value_color_1,
	shift_value_color_2, shift_value_color_3, Metric_Number uint8
	job_limit_type, shift_limit_type string
}

//////////////////////////////////////////////////////////////////////////////////////
//                                                                                  //
func MSSQL_get_device_setting(altix_device_id int64) ([]byte, []byte, uint32, error) {
	//                                                                                  //
	//////////////////////////////////////////////////////////////////////////////////////
	data := make([]byte, 2 + 208) // khusus yg ini, + 2 dari data. yg lain + 3 arena ada recCount()
	data[0] = 0
	data[1] = 7

	qry :=
		// uint32
		"select machine_Id, convert(char(24),machine_Name), " +
		"  convert(int, when_no_job_Scale_Total_Count * 1000.0), " +
		"  convert(int, when_no_job_Scale_Reject_Count * 1000.0), " +

		"  Screen_Duration_mSec1, Screen_Duration_mSec2, " +
		"  when_no_job_Down_Rate_per_Hour, when_no_job_Run_Rate_per_Hour, " +
		"  when_no_job_Slow_Cycle_Qty_Sequence_Display_Treshold, " +

		"  isnull(Screen1Row1_Metric_Id,0), isnull(Screen1Row2_Metric_Id,0), " +
		"  isnull(Screen1Row3_Metric_Id,0), isnull(Screen1Row4_Metric_Id,0), " +
		"  isnull(Screen2Row1_Metric_Id,0), isnull(Screen2Row2_Metric_Id,0), " +
		"  isnull(Screen2Row3_Metric_Id,0), isnull(Screen2Row4_Metric_Id,0), " +

		// uint8
		"  Screen1Row1_JobShift, Screen1Row2_JobShift, Screen1Row3_JobShift, Screen1Row4_JobShift, " +
		"  Screen2Row1_JobShift, Screen2Row2_JobShift, Screen2Row3_JobShift, Screen2Row4_JobShift, " +

		"  Brightness, isnull(Altix_Device_TotalCount_Gang,1), isnull(Altix_Device_NotGood_Gang,2), " +
		"  isnull(mr.RemoteSetting_Number,0), convert(int, Show_Active_Job), " +
		"  convert(int, Show_Shift_Info), convert(int, Target_Counter_Stop_while_Down), " +
		"  convert(int, Target_Counter_Stop_while_Standby), convert(int, Target_Counter_Stop_while_Setup), " +

		"  when_no_job_Switch_from_Run_to_Down_based_on, when_no_job_Switch_from_Down_to_Run_based_on, " +
		"  when_no_job_Switch_from_Down_to_Run_Count, " +
		//"  when_no_job_Switch_from_Down_to_Run_counter_of, " +
		"  when_no_job_Rate_per_Hour_count_from, " +

		"  convert(char(19),Time_Stamp,120), " +
		"  isnull(convert(char(19),Next_TimeSchedule_StartDate,120), '1900-01-01 00:00:00'), " +
		"  isnull(convert(char(19),Occasional_TimeSchedule_StartDate,120), '1900-01-01 00:00:00'), " +
		"  isnull(convert(char(19),Occasional_TimeSchedule_EndDate,120), '1900-01-01 00:00:00'), " +

		"  convert(bigint, (case when isnumeric(when_no_job_Ideal_Cycle_Pieces) = 0 then 0 else when_no_job_Ideal_Cycle_Pieces end) * 1000000.0), " +
		"  convert(bigint, (case when isnumeric(when_no_job_Takt_Pieces) = 0 then 0 else when_no_job_Takt_Pieces end) * 1000000.0), " +

		"  convert(int, (case when isnumeric(when_no_job_Ideal_Cycle_Hours) = 0 then 0 else when_no_job_Ideal_Cycle_Hours end) * 1000000.0), " +
		"  convert(int, (case when isnumeric(when_no_job_Ideal_Cycle_Minutes) = 0 then 0 else when_no_job_Ideal_Cycle_Minutes end) * 1000000.0), " +
		"  convert(int, (case when isnumeric(when_no_job_Ideal_Cycle_Seconds) = 0 then 0 else when_no_job_Ideal_Cycle_Seconds end) * 1000000.0), " +
		"  convert(int, (case when isnumeric(when_no_job_Takt_Hours) = 0 then 0 else when_no_job_Takt_Hours end) * 1000000.0), " +
		"  convert(int, (case when isnumeric(when_no_job_Takt_Minutes) = 0 then 0 else when_no_job_Takt_Minutes end) * 1000000.0), " +
		"  convert(int, (case when isnumeric(when_no_job_Takt_Seconds) = 0 then 0 else when_no_job_Takt_Seconds end) * 1000000.0), " +

		"  when_no_job_Slow_Cycle_Pct, when_no_job_Small_Stop_Pct, when_no_job_Full_Down_Pct, " +

		"  Debounce_Gang_1_mSec, Debounce_Gang_2_mSec, Debounce_Gang_3_mSec, Debounce_Gang_4_mSec " +
		"from mstMachine as mm " +
		"left join mstRemoteSetting as mr on mr.RemoteSetting_Id = mm.RemoteSetting_Id " +

		"where Altix_Device_Id = " + strconv.FormatInt(altix_device_id, 10) + " "
	rows, err := db.Query(qry)
	defer rows.Close()

	deviceTimeStamp := make([]byte, 7)
	if err != nil {
		fmt.Println("MSSQL_get_device_setting err:\n" + err.Error() + "\n" + qry)
		return data, deviceTimeStamp, 0, err
	}

	if rows.Next() {
	} else {
		fmt.Println("MSSQL_get_device_setting err: 0 rows.\n" + qry)
		return data, deviceTimeStamp, 0, errors.New("MSSQL_get_device_setting err: 0 rows.\n" + qry)
	}

	var machine_Id, jumlahId    uint32
	var machine_name       string
	var	when_no_job__scale_total_count, 
		when_no_job__scale_reject_count, 
		Screen_Duration_1_mSec, Screen_Duration_2_mSec,
		when_no_job__down__rate_per_hour, 
		when_no_job__run__rate_per_hour,
		when_no_job__slow_cycle__treshold_sequence_display uint32

	metricId := make([]uint32, 8)

	job_shift := make([]string, 8)

	var brightness, total_count_gang, reject_count_gang,
		remote_setting, show_active_job, show_shift_info,
		target_counter_stopped_while_down, 
		target_counter_stopped_while_standbby,
		target_counter_stopped_while_setup,
		when_no_job__switch_from_run_to_down__based_on,
		when_no_job__switch_from_down_to_run__based_on,
		when_no_job__switch_from_down_to_run__count,
		//when_no_job__switch_from_down_to_run__counter_of,
		when_no_job__rate_per_hour__based_on__counter_of  int

	var device_setting_timeStamp, next_timeSchedule_startDate, 
		occas_timeSchedule_startDate, occas_timeSchedule_endDate string 

	var when_no_job__ideal_cycle_pieces_x_1M, when_no_job__takt_pieces_x_1M uint64

	var when_no_job__ideal_cycle_hours_x_1M, 
		when_no_job__ideal_cycle_minutes_x_1M,
		when_no_job__ideal_cycle_seconds_x_1M,
		when_no_job__takt_hours_x_1M, when_no_job__takt_minutes_x_1M,
		when_no_job__takt_seconds_x_1M uint32

	var when_no_job__pct_slow_cycle, when_no_job__pct_small_stop, 
	    when_no_job__pct_full_stop uint16
		
	gangDebounce_mSec := make([]uint16, 4)

	//////////////
	// ROW SCAN //
	//////////////
	err = rows.Scan(&machine_Id, &machine_name,
		&when_no_job__scale_total_count, &when_no_job__scale_reject_count,

		&Screen_Duration_1_mSec, &Screen_Duration_2_mSec,
		&when_no_job__down__rate_per_hour, &when_no_job__run__rate_per_hour,
		&when_no_job__slow_cycle__treshold_sequence_display,

		&metricId[0], &metricId[1], &metricId[2], &metricId[3],
		&metricId[4], &metricId[5], &metricId[6], &metricId[7],

		&job_shift[0], &job_shift[1], &job_shift[2], &job_shift[3],
		&job_shift[4], &job_shift[5], &job_shift[6], &job_shift[7],

		&brightness, &total_count_gang, &reject_count_gang,
		&remote_setting, &show_active_job, &show_shift_info,
		&target_counter_stopped_while_down, &target_counter_stopped_while_standbby,
		&target_counter_stopped_while_setup,

		&when_no_job__switch_from_run_to_down__based_on,
		&when_no_job__switch_from_down_to_run__based_on,
		//&when_no_job__switch_from_down_to_run__counter_of, 
		&when_no_job__switch_from_down_to_run__count,
		&when_no_job__rate_per_hour__based_on__counter_of,

		&device_setting_timeStamp, &next_timeSchedule_startDate,
		&occas_timeSchedule_startDate, &occas_timeSchedule_endDate,

		&when_no_job__ideal_cycle_pieces_x_1M, &when_no_job__takt_pieces_x_1M,

		&when_no_job__ideal_cycle_hours_x_1M, &when_no_job__ideal_cycle_minutes_x_1M,
		&when_no_job__ideal_cycle_seconds_x_1M, 
		&when_no_job__takt_hours_x_1M, &when_no_job__takt_minutes_x_1M,
		&when_no_job__takt_seconds_x_1M,

		&when_no_job__pct_slow_cycle, &when_no_job__pct_small_stop,
		&when_no_job__pct_full_stop,

		&gangDebounce_mSec[0], &gangDebounce_mSec[1], &gangDebounce_mSec[2], &gangDebounce_mSec[3])
	if err != nil {
		fmt.Println("MSSQL_get_device_setting err: rows.Scan\n" + err.Error())
		return data, deviceTimeStamp, 0, err
	}

	// validasi dan koreksi di sini agar tidak memberatkankerja MCU
	if when_no_job__slow_cycle__treshold_sequence_display > 100 { when_no_job__slow_cycle__treshold_sequence_display = 100 }
	if brightness < 0 ||brightness > 10 { brightness = 5 }
	if total_count_gang < 1 || total_count_gang > 4 { total_count_gang = 1 }
	if reject_count_gang < 1 || reject_count_gang > 4 { reject_count_gang = 2 }
	if remote_setting < 0 || remote_setting > 6 { remote_setting = 0 }
	if show_active_job < 0 || show_active_job > 1 { show_active_job = 0 }
	if show_shift_info < 0 || show_shift_info > 1 { show_shift_info = 0 }
	if target_counter_stopped_while_down < 0 || target_counter_stopped_while_down > 1 { target_counter_stopped_while_down = 0 }
	if target_counter_stopped_while_standbby < 0 || target_counter_stopped_while_standbby > 1 { target_counter_stopped_while_standbby = 1 }
	if target_counter_stopped_while_setup < 0 || target_counter_stopped_while_setup > 1 { target_counter_stopped_while_setup = 1 }

	if when_no_job__switch_from_run_to_down__based_on < 0 || when_no_job__switch_from_run_to_down__based_on > 1 {
		when_no_job__switch_from_run_to_down__based_on = 0
	}
	if when_no_job__switch_from_down_to_run__based_on < 0 || when_no_job__switch_from_down_to_run__based_on > 2 {
		when_no_job__switch_from_down_to_run__based_on = 0
	}
	if when_no_job__switch_from_down_to_run__count < 0 || when_no_job__switch_from_down_to_run__count > 100 {
		when_no_job__switch_from_down_to_run__count = 10
	}
	if when_no_job__rate_per_hour__based_on__counter_of < 0 || when_no_job__rate_per_hour__based_on__counter_of > 1 {
		when_no_job__rate_per_hour__based_on__counter_of = 0
	}

	nc_write_uint32_to_byte_slice__LSB_to_MSB(machine_Id, &data, 2)

	mnX := np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(machine_name,24)
	for s := 0; s < 24; s++ {
		data[6 + s] = mnX[s]
	}
	data[29] = 0 // karena byte terakhir string di c/c++ harus chr(0)

	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__scale_total_count, &data, 30)
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__scale_reject_count, &data, 34)

	nc_write_uint32_to_byte_slice__LSB_to_MSB(Screen_Duration_1_mSec, &data, 38)
	nc_write_uint32_to_byte_slice__LSB_to_MSB(Screen_Duration_2_mSec, &data, 42)
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__down__rate_per_hour, &data, 46)
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__run__rate_per_hour, &data, 50)

	nc_write_uint32_to_byte_slice__LSB_to_MSB(
		when_no_job__slow_cycle__treshold_sequence_display, &data, 54)

	jumlahId = 0;
	for i := 0; i < 8; i++ {
		jumlahId += metricId[i]
		nc_write_uint32_to_byte_slice__LSB_to_MSB(metricId[i], &data, 58+i*4)
		if job_shift[i][0] == 'J' {	data[i+90] = 'J' }
		if job_shift[i][0] == 'S' {	data[i+90] = 'S' }
	}

	ndx := 98
	data[ndx] = byte(brightness); ndx++
	data[ndx] = byte(total_count_gang); ndx++
	data[ndx] = byte(reject_count_gang); ndx++
	data[ndx] = byte(remote_setting); ndx++
	data[ndx] = byte(show_active_job); ndx++
	data[ndx] = byte(show_shift_info); ndx++
	data[ndx] = byte(target_counter_stopped_while_down); ndx++
	data[ndx] = byte(target_counter_stopped_while_standbby); ndx++
	data[ndx] = byte(target_counter_stopped_while_setup); ndx++

	data[ndx] = byte(when_no_job__switch_from_run_to_down__based_on); ndx++
	data[ndx] = byte(when_no_job__switch_from_down_to_run__based_on); ndx++
	data[ndx] = byte(when_no_job__switch_from_down_to_run__count); ndx++
	//data[ndx] = byte(when_no_job__switch_from_down_to_run__counter_of); ndx++
	data[ndx] = byte(when_no_job__rate_per_hour__based_on__counter_of); ndx++

	//fmt.Println("before timeStamp, ndx=",ndx)
	np_write_datetime120_to_byte_slice7(device_setting_timeStamp, &data, ndx); ndx += 7
	for i := 0; i < 7; i++ { deviceTimeStamp[i] = data[ndx-7+i] }

	np_write_datetime120_to_byte_slice7(next_timeSchedule_startDate, &data, ndx); ndx += 7
	np_write_datetime120_to_byte_slice7(occas_timeSchedule_startDate, &data, ndx); ndx += 7
	np_write_datetime120_to_byte_slice7(occas_timeSchedule_endDate, &data, ndx); ndx += 7

	fmt.Println("before lompat 7, ndx=",ndx)
	ndx += 7 // dummy

	// 64 bit 
	nc_write_uint64_to_byte_slice__LSB_to_MSB(when_no_job__ideal_cycle_pieces_x_1M, &data, ndx); ndx += 8
	nc_write_uint64_to_byte_slice__LSB_to_MSB(when_no_job__takt_pieces_x_1M, &data, ndx); ndx += 8

	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__ideal_cycle_hours_x_1M, &data, ndx); ndx += 4
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__ideal_cycle_minutes_x_1M, &data, ndx); ndx += 4
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__ideal_cycle_seconds_x_1M, &data, ndx); ndx += 4
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__takt_hours_x_1M, &data, ndx); ndx += 4
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__takt_minutes_x_1M, &data, ndx); ndx += 4
	nc_write_uint32_to_byte_slice__LSB_to_MSB(when_no_job__takt_seconds_x_1M, &data, ndx); ndx += 4

	nc_write_uint16_to_byte_slice__LSB_to_MSB(when_no_job__pct_slow_cycle, &data, ndx); ndx += 2
	nc_write_uint16_to_byte_slice__LSB_to_MSB(when_no_job__pct_small_stop, &data, ndx); ndx += 2
	nc_write_uint16_to_byte_slice__LSB_to_MSB(when_no_job__pct_full_stop, &data, ndx); ndx += 2

	//fmt.Println("ndx=",ndx)
	for i := 0; i < 4; i++ {
		nc_write_uint16_to_byte_slice__LSB_to_MSB(gangDebounce_mSec[i], &data, ndx + i*2); ndx += 2
	}

	ndx += 2 // dummy

	nc_write_uint32_to_byte_slice__LSB_to_MSB(0x44564353, &data, ndx)
	//fmt.Println("ndx=",ndx)

	return data, deviceTimeStamp, jumlahId, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                 //
func MSSQL_get_mstMetrics(tenant_id int64, last_timeStamp string) (
	[]uint32, []altix_mstMetric, error, error) {
	//                                                                                                 //
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	var array_Id []uint32
	var array_mstMetric []altix_mstMetric

	qry := "set dateformat ymd " +
		"select Metric_Id from mstMetric where Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"order by Metric_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_mstMetrics error:\n" + listErr.Error() + "\n" + qry)
		return array_Id, array_mstMetric, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return array_Id, array_mstMetric, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1390 - 4 - 4*len(array_Id)
	// asumsi pjgData mstMetric = 60
	maxRec := jSisaData / 60

	qry = "set dateformat ymd " +
		"select top " + strconv.FormatInt(int64(maxRec), 10) +
		"  Metric_Id, convert(char(15),Metric_Desc), convert(char(15),Display_As), " +
		"  Metric_Number, JobLabel_Color, ShiftLabel_Color, " +
		"  convert(char(1),Job_Limit_Type), convert(char(1),Shift_Limit_Type), " +
		"  convert(int, JobLimit1 * 1000.0), convert(int, JobLimit2 * 1000.0), " +
		"  convert(int, ShiftLimit1 * 1000.0), convert(int, ShiftLimit2 * 1000.0), " +
		"  JobValue_Color1, JobValue_Color2, JobValue_Color3, ShiftValue_Color1, " +
		"  ShiftValue_Color2, ShiftValue_Color3, convert(varchar(19),Time_Stamp,120) " +
		"from mstMetric where Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" +
		last_timeStamp + "') " +
		"order by Time_Stamp"
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_mstMetrics error:\n" + err.Error() + "\n" + qry)
		return array_Id, array_mstMetric, nil, err
	}
	for rows.Next() {
		mm := altix_mstMetric{}
		if err := rows.Scan(&mm.metric_Id, &mm.metric_Desc, &mm.display_As, &mm.Metric_Number,
			&mm.job_label_color, &mm.shift_label_color, &mm.job_limit_type, &mm.shift_limit_type,
			&mm.job_limit_1_x_1K, &mm.job_limit_2_x_1K, &mm.shift_limit_1_x_1K, &mm.shift_limit_2_x_1K,
			&mm.job_value_color_1, &mm.job_value_color_2, &mm.job_value_color_3, &mm.shift_value_color_1,
			&mm.shift_value_color_2, &mm.shift_value_color_3, &mm.str_time_stamp); err != nil {
			return array_Id, array_mstMetric, nil, err
		}
		array_mstMetric = append(array_mstMetric, mm)
		fmt.Println(mm)
	}
	return array_Id, array_mstMetric, nil, nil
}