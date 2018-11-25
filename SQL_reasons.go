package main

import (
	//"time"
	//"database/sql"
	//"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
	//"strings"
)

type altix_mstReason struct {
	Reason_Id                   uint32
	Reason_Desc, str_time_stamp string
	Remote_Code                 uint8
	Green_Duration, Yellow_Duration, Reason_Duration,
	GoodCount_to_End_Setup uint16
	Is_Deleted uint8
}

//////////////////////////////////////////////////////////////////////
//                                                                  //
func MSSQL_get_mstReasons(Reason_Type string, tenant_id int64,
	last_timeStamp string) (lastDelete_dateTime string, array_Id []uint32,
	array_mstReason []altix_mstReason, list error, data error) {
	//                                                                  //
	//////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_mstReason []altix_mstReason

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select Reason_Id from mstReason " +
		"where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"  and Reason_Type = '" + Reason_Type + "' " +
		"order by Reason_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_mstReasons error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_mstReason, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_mstReason, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1380 - 4 - 4*len(array_Id)

	maxRec := 0
	if Reason_Type == "Down" {
		maxRec = jSisaData / 48
	}
	if Reason_Type == "Standby" {
		maxRec = jSisaData / 52
	}
	if Reason_Type == "Setup" {
		maxRec = jSisaData / 56
	}

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("mstReason              ", Reason_Type, tenant_id)
	if errD != nil {
		return "", array_Id, array_mstReason, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		"select  top " + strconv.FormatInt(int64(maxRec), 10) +
		"  Reason_Id, convert(varchar(19),Time_Stamp,120), " +
		"  Remote_Code, convert(char(29),Reason_Desc), Green_Duration_In_Sec," +
		"  Yellow_Duration_In_Sec, Reason_Duration_In_Second, GoodCount_to_End_Setup, " +
		"  convert(int,Is_Deleted) " +
		"from mstReason where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"  and Reason_Type = '" + Reason_Type + "' order by Is_Deleted desc, Time_Stamp"
	//fmt.Println(qry)
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_mstReasons " + Reason_Type + " error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_mstReason, nil, err
	}
	for rows.Next() {
		mm := altix_mstReason{}
		if err := rows.Scan(&mm.Reason_Id, &mm.str_time_stamp, &mm.Remote_Code,
			&mm.Reason_Desc, &mm.Green_Duration, &mm.Yellow_Duration, &mm.Reason_Duration,
			&mm.GoodCount_to_End_Setup, &mm.Is_Deleted); err != nil {
			return "", array_Id, array_mstReason, nil, err
		}
		array_mstReason = append(array_mstReason, mm)
		fmt.Println(mm)
	}
	return qryD[1],array_Id, array_mstReason, nil, nil
}