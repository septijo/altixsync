package main

import (
	//"time"
	//"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
	//"strings"
)

type altix_timeSchedule_Shift struct {
	TimeSchedule_Shift_Id, Shift_Id      uint32
	Shift_Minute_Duration                uint16
	Shift_Start_HhMm, str_time_stamp     string
	Shift_Start_DayOfWeek, CurrentOrNext uint8
	Is_Deleted                           uint8
}

type altix_timeSchedule_ShiftBreak struct {
	TimeSchedule_ShiftBreak_Id, TimeSchedule_Shift_Id, Break_Id uint32
	Break_Minute_Duration                                       uint16
	Break_Start_HhMm, str_time_stamp                            string
	Break_Start_DayOfWeek, Current_Next_Occasional              uint8
	Is_Deleted                                                  uint8
}

////////////////////////////////////////////////////////////////////////////////////
//                                                                                //
func MSSQL_get_timeSchedule_Shifts(machine_id int64, last_timeStamp string) (
	lastDelete_dateTime string, array_Id []uint32, array_timeSchedule_Shift []altix_timeSchedule_Shift, list error, data error) {
	//                                                                                //
	////////////////////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_timeSchedule_Shift []altix_timeSchedule_Shift

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select TimeSchedule_Shift_Id from TimeSchedule_Shift " +
		"where Is_Deleted = 0 and Machine_Id = " + strconv.FormatInt(machine_id, 10) + " " +
		"order by TimeSchedule_Shift_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_timeSchedule_Shifts error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_timeSchedule_Shift, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_timeSchedule_Shift, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1380 - 4 - 4*len(array_Id)
	// asumsi pjgData TimeSchedule_Shift = 28
	maxRec := jSisaData / 28

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("TimeSchedule_Shift", "", machine_id)
	if errD != nil {
		return "", array_Id, array_timeSchedule_Shift, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		"select top " + strconv.FormatInt(int64(maxRec), 10) +
		"  TimeSchedule_Shift_Id, Shift_Id, Shift_Minute_Duration, " +
		"  convert(varchar(19),Time_Stamp,120), Shift_Start_HhMm, Shift_Start_DayOfWeek, " +
		"  CurrentOrNext, convert(int,Is_Deleted) " +
		"from TimeSchedule_Shift where Is_Deleted = 0 and Machine_Id = " + strconv.FormatInt(machine_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"order by Is_Deleted desc, Time_Stamp"
	//fmt.Println(qry)
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_timeSchedule_Shifts error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_timeSchedule_Shift, nil, err
	}
	for rows.Next() {
		mm := altix_timeSchedule_Shift{}
		if err := rows.Scan(&mm.TimeSchedule_Shift_Id, &mm.Shift_Id, &mm.Shift_Minute_Duration,
			&mm.str_time_stamp, &mm.Shift_Start_HhMm, &mm.Shift_Start_DayOfWeek, &mm.CurrentOrNext,
			&mm.Is_Deleted); err != nil {
			return "", array_Id, array_timeSchedule_Shift, nil, err
		}
		array_timeSchedule_Shift = append(array_timeSchedule_Shift, mm)
		fmt.Println(mm)
	}
	return qryD[1], array_Id, array_timeSchedule_Shift, nil, nil
}

////////////////////////////////////////////////////////////////////////////////////
//                                                                                //
func MSSQL_get_timeSchedule_ShiftBreaks(machine_id int64, last_timeStamp string) (lastDelete_dateTime string, 
	array_Id []uint32, array_timeSchedule_ShiftBreak []altix_timeSchedule_ShiftBreak, list error, data error) {
	//                                                                                //
	////////////////////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_timeSchedule_ShiftBreak []altix_timeSchedule_ShiftBreak

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select TSSB.TimeSchedule_ShiftBreak_Id " +
		"from TimeSchedule_ShiftBreak as TSSB " +
		"left join TimeSchedule_Shift as TSS on TSS.TimeSchedule_Shift_Id = TSSB.TimeSchedule_Shift_Id " +
		"where TSSB.Is_Deleted = 0 and TSS.Machine_Id = " + strconv.FormatInt(machine_id, 10) + " " +
		"order by TSSB.TimeSchedule_ShiftBreak_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_timeSchedule_ShiftBreaks error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_timeSchedule_ShiftBreak, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_timeSchedule_ShiftBreak, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	//jSisaData := 1390 - 4 - 4 * len(array_Id)
	// asumsi pjgData TSSB = 32
	//maxRec := jSisaData / 32

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("TimeSchedule_Shift", "", machine_id)
	if errD != nil {
		return "", array_Id, array_timeSchedule_ShiftBreak, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		"select " + //top " + strconv.FormatInt(int64(maxRec), 10) +
		"  TS_SB.TimeSchedule_ShiftBreak_Id, TS_SB.TimeSchedule_Shift_Id, " +
		"  TS_SB.Break_Id, TS_SB.Break_Dur, convert(varchar(19),TS_SB.Time_Stamp,120), " +
		"  TS_SB.Break_Start, TS_SB.Break_Start_DayOfWeek, TS_SB.CurrentOrNext, " +
		"  convert(int,TS_SB.Is_Deleted) " +
		"from TimeSchedule_ShiftBreak as TS_SB " +
		"left join TimeSchedule_Shift as TS_S on TS_S.TimeSchedule_Shift_Id = " +
		"  TS_SB.TimeSchedule_Shift_Id " +
		"where TS_SB.Is_Deleted = 0 and TS_S.Machine_Id = " + strconv.FormatInt(machine_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),TS_SB.Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"order by TS_SB.Is_Deleted desc, TS_SB.Time_Stamp"
	//fmt.Println(qry)
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_timeSchedule_ShiftBreaks error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_timeSchedule_ShiftBreak, nil, err
	}
	for rows.Next() {
		mm := altix_timeSchedule_ShiftBreak{}
		if err := rows.Scan(&mm.TimeSchedule_ShiftBreak_Id, &mm.TimeSchedule_Shift_Id,
			&mm.Break_Id, &mm.Break_Minute_Duration, &mm.str_time_stamp,
			&mm.Break_Start_HhMm, &mm.Break_Start_DayOfWeek, &mm.Current_Next_Occasional,
			&mm.Is_Deleted); err != nil {
			return "", array_Id, array_timeSchedule_ShiftBreak, nil, err
		}
		array_timeSchedule_ShiftBreak = append(array_timeSchedule_ShiftBreak, mm)
		fmt.Println(mm)
	}
	return qryD[1], array_Id, array_timeSchedule_ShiftBreak, nil, nil
}