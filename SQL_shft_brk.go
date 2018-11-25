package main

import (
	//"time"
	//"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
	//"strings"
)

type altix_mstShift struct {
	Shift_Id            uint32
	str_time_stamp      string
	Shift_Number, Color uint8
	Is_Deleted          uint8
}

type altix_mstBreak struct {
	Break_Id                   uint32
	Break_Desc, str_time_stamp string
	Color                      uint8
	Is_Deleted                 uint8
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                //
func MSSQL_get_mstShifts(tenant_id int64, last_timeStamp string) (lastDelete_dateTime string, 
	array_Id []uint32, array_mstShift []altix_mstShift, list error, data error) {
	//                                                                                                //
	////////////////////////////////////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_mstShift []altix_mstShift

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select Shift_Id from mstShift " +
		"where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"order by Shift_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_mstShifts error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_mstShift, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_mstShift, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1380 - 4 - 4*len(array_Id)
	// asumsi pjgData mstShift = 28
	maxRec := jSisaData / 28

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("mstShift", "", tenant_id)
	if errD != nil {
		return "", array_Id, array_mstShift, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		"select top " + strconv.FormatInt(int64(maxRec), 10) +
		"  Shift_Id, convert(varchar(19),Time_Stamp,120), Shift_Number, Color, " +
		"  convert(int,Is_Deleted) " +
		"from mstShift where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"order by Is_Deleted desc, Time_Stamp"
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_mstShifts error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_mstShift, nil, err
	}
	for rows.Next() {
		mm := altix_mstShift{}
		if err := rows.Scan(&mm.Shift_Id, &mm.str_time_stamp, &mm.Shift_Number, &mm.Color,
			&mm.Is_Deleted); err != nil {
			return "", array_Id, array_mstShift, nil, err
		}
		array_mstShift = append(array_mstShift, mm)
		fmt.Println(mm)
	}
	return qryD[1], array_Id, array_mstShift, nil, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                //
func MSSQL_get_mstBreaks(tenant_id int64, last_timeStamp string) (lastDelete_dateTime string, 
	array_Id []uint32, array_mstBreak []altix_mstBreak, list error, data error) {
	//                                                                                                //
	////////////////////////////////////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_mstBreak []altix_mstBreak

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select Break_Id from mstBreak " +
		"where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"order by Break_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_mstBreaks error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_mstBreak, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_mstBreak, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1380 - 4 - 4*len(array_Id)
	// asumsi pjgData mstBreak = 48
	maxRec := jSisaData / 48

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("mstBreak", "", tenant_id)
	if errD != nil {
		return "", array_Id, array_mstBreak, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		"select top " + strconv.FormatInt(int64(maxRec), 10) +
		"  Break_Id, convert(varchar(19), Time_Stamp,120), Color, " +
		"  convert(char(29),Break_Desc), convert(int,Is_Deleted) " +
		"from mstBreak where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"order by Is_Deleted desc, Time_Stamp"
	fmt.Println(qry)
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_mstBreaks error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_mstBreak, nil, err
	}
	for rows.Next() {
		mm := altix_mstBreak{}
		if err := rows.Scan(&mm.Break_Id, &mm.str_time_stamp, &mm.Color, &mm.Break_Desc,
			&mm.Is_Deleted); err != nil {
			return "", array_Id, array_mstBreak, nil, err
		}
		array_mstBreak = append(array_mstBreak, mm)
		fmt.Println(mm)
	}
	return qryD[1], array_Id, array_mstBreak, nil, nil
}
