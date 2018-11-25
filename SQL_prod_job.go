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

type altix_mstProduct struct {
	Ideal_Cycle_Pieces_x_1M, Takt_Pieces_x_1M uint64
	Product_Id, Ideal_Cycle_Hours_x_1M, Ideal_Cycle_Minutes_x_1M,
	Ideal_Cycle_Seconds_x_1M, Takt_Hours_x_1M, Takt_Minutes_x_1M,
	Takt_Seconds_x_1M, Slow_Cycle__Treshold_Sequence_Display,
	Scale_Total_Count_x_1K, Scale_Reject_Count_x_1K uint32
	str_time_stamp, Product_Name                  string
	Pct_Slow_Cycle, Pct_Small_Stop, Pct_Full_Stop uint16
	Is_Deleted                                    uint8
}

type altix_job struct {
	Ideal_Cycle_Pieces_x_1M, Takt_Pieces_x_1M uint64
	Job_Id, Product_Id, Ideal_Cycle_Hours_x_1M,
	Ideal_Cycle_Minutes_x_1M, Ideal_Cycle_Seconds_x_1M,
	Takt_Hours_x_1M, Takt_Minutes_x_1M, Takt_Seconds_x_1M,
	Scale_Total_Count_x_1K, Scale_Reject_Count_x_1K,
	Goal_Qty, Slow_Cycle__Treshold_Sequence_Display,
	Down_Rate_per_Hour, Run_Rate_per_Hour uint32
	str_time_stamp, Job_Desc, Ref_No              string
	Pct_Slow_Cycle, Pct_Small_Stop, Pct_Full_Stop uint16
	Switch_from_Run_to_Down_based_on, Switch_from_Down_to_Run_based_on, 
	Switch_from_Down_to_Run_Count, //Switch_from_Down_to_Run_counter_of, 
	Rate_per_Hour_counter_of, Remote_Code, Is_Deleted  uint8
}
////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                //
func MSSQL_get_mstProducts(tenant_id int64, last_timeStamp string) (lastDelete_dateTime string, 
	array_Id []uint32, array_mstProduct []altix_mstProduct, list error, data error) {
	//                                                                                                //
	////////////////////////////////////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_mstProduct []altix_mstProduct

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select Product_Id from mstProduct " +
		"where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"order by Product_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_mstProducts error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_mstProduct, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_mstProduct, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1380 - 4 - 4*len(array_Id)
	// asumsi pjgData mstProduct = 104
	maxRec := jSisaData / 104

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("mstProduct", "", tenant_id)
	if errD != nil {
		return "", array_Id, array_mstProduct, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		// 84 bytes jadi 1397 / 84 = 16
		"select top " + strconv.FormatInt(int64(maxRec), 10) +
		"  convert(bigint, Ideal_Cycle_Pieces * 1000000.0), convert(bigint, Takt_Pieces * 1000000.0), " +

		"  Product_Id, convert(varchar(19), Time_Stamp,120), " +

		"  Slow_Cycle_Pct, Small_Stop_Pct, Full_Down_Pct, " +

		"  convert(char(29),Product_Name), convert(int, Ideal_Cycle_Hours * 1000000.0), " +
		"  convert(int, Ideal_Cycle_Minutes * 1000000.0), convert(int, Ideal_Cycle_Seconds * 1000000.0), " +

		"  convert(int, Takt_Hours * 1000000.0), convert(int, Takt_Minutes * 1000000.0), convert(int, Takt_Seconds * 1000000.0), " +

		"  convert(int, Scale_Total_Count * 1000.0), convert(int, Scale_Reject_Count * 1000.0), " +
		"  Treshold_Sequence_Display as Slow_Cycle_Treshold_Sequence_Display, " +
		"  convert(int,Is_Deleted) " +
		"from mstProduct where Is_Deleted = 0 and Tenant_Id = " + strconv.FormatInt(tenant_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"order by Is_Deleted desc, Time_Stamp"
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_mstProducts error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_mstProduct, nil, err
	}
	for rows.Next() {
		mm := altix_mstProduct{}
		//fmt.Println("before scan mstProducts")
		if err := rows.Scan(&mm.Ideal_Cycle_Pieces_x_1M, &mm.Takt_Pieces_x_1M,

			&mm.Product_Id, &mm.str_time_stamp,

			&mm.Pct_Slow_Cycle, &mm.Pct_Small_Stop, &mm.Pct_Full_Stop,

			&mm.Product_Name, &mm.Ideal_Cycle_Hours_x_1M,
			&mm.Ideal_Cycle_Minutes_x_1M, &mm.Ideal_Cycle_Seconds_x_1M,

			&mm.Takt_Hours_x_1M, &mm.Takt_Minutes_x_1M, &mm.Takt_Seconds_x_1M,

			&mm.Scale_Total_Count_x_1K, &mm.Scale_Reject_Count_x_1K,

			&mm.Slow_Cycle__Treshold_Sequence_Display,

			&mm.Is_Deleted); err != nil {
			return "", array_Id, array_mstProduct, nil, err
		}
		//fmt.Println("after scan mstProducts")
		array_mstProduct = append(array_mstProduct, mm)
		fmt.Println(mm)
	}
	return qryD[1], array_Id, array_mstProduct, nil, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                //
func MSSQL_get_jobs(machine_id int64, last_timeStamp string) (lastDelete_dateTime string, 
	array_Id []uint32, array_job []altix_job, list error, data error) {
	//                                                                                                //
	////////////////////////////////////////////////////////////////////////////////////////////////////
	//var array_Id []uint32
	//var array_job []altix_job

	// step 1 : ambil List Id
	qry := "set dateformat ymd " +
		"select Job_Id from Job " +
		"where Is_Deleted = 0 and Machine_Id = " + strconv.FormatInt(machine_id, 10) + " " +
		"order by Job_Id"
	listRows, listErr := db.Query(qry)
	defer listRows.Close()

	if listErr != nil {
		fmt.Println("MSSQL_get_mstMetrics error:\n" + listErr.Error() + "\n" + qry)
		return "", array_Id, array_job, listErr, nil
	}
	for listRows.Next() {
		var mId uint32
		if listErr := listRows.Scan(&mId); listErr != nil {
			return "", array_Id, array_job, listErr, nil
		}
		array_Id = append(array_Id, mId)
	}
	fmt.Println(array_Id)

	// hitung max rec agar muat di 1400 MTU
	jSisaData := 1380 - 4 - 4*len(array_Id)
	// asumsi pjgData mstMetric = 140
	maxRec := jSisaData / 144

	// step 2 : ambil Last Tanggal Delete
	qryD, errD := get_LastTimeStamp("job", "", machine_id)
	if errD != nil {
		return "", array_Id, array_job, nil, errD
	}

	// step 3 : ambil data nya
	qry = "set dateformat ymd " +
		// 128 bytes jadi 1397 / 84 = 10
		"select top " + strconv.FormatInt(int64(maxRec), 10) +
		"  convert(bigint, Ideal_Cycle_Pieces * 1000000.0), convert(bigint, Takt_Pieces * 1000000.0), " +

		"  Job_Id, Product_Id, convert(varchar(19), Time_Stamp,120), Slow_Cycle_Pct, " +
		"  Small_Stop_Pct, Full_Down_Pct, convert(char(29),Job_Desc), convert(char(29),Ref_No), " +
		"  convert(int, Ideal_Cycle_Hours * 1000000.0), convert(int, Ideal_Cycle_Minutes * 1000000.0), " +
		"  convert(int, Ideal_Cycle_Seconds * 1000000.0), " +
		"  convert(int, Takt_Hours * 1000000.0), convert(int, Takt_Minutes * 1000000.0), " +
		"  convert(int, Takt_Seconds * 1000000.0), " +
		"  convert(int, Scale_Total_Count * 1000.0), convert(int, Scale_Reject_Count * 1000.0), " +
		"  Goal_Qty, Slow_Cycle_Qty_Sequence_Display_Treshold, Down_Rate_per_Hour, Run_Rate_per_Hour, " +
		"  Switch_from_Run_to_Down_based_on, Switch_from_Down_to_Run_based_on, " +
		"  Switch_from_Down_to_Run_Count, Rate_per_Hour_count_from, " + //Switch_from_Down_to_Run_counter_of, 
		"  Remote_Code, convert(int,Is_Deleted) " +
		"from Job where Is_Deleted = 0 and Machine_Id = " + strconv.FormatInt(machine_id, 10) + " " +
		"  and convert(DateTime,convert(char(19),Time_Stamp,120)) > convert(DateTime,'" + last_timeStamp + "') " +
		"order by Is_Deleted desc, Time_Stamp"
	rows, err := db.Query(qry)
	defer rows.Close()
	if err != nil {
		fmt.Println("MSSQL_get_jobs error:\n" + err.Error() + "\n" + qry)
		return "", array_Id, array_job, nil, err
	}
	for rows.Next() {
		mm := altix_job{}
		if err := rows.Scan(&mm.Ideal_Cycle_Pieces_x_1M, &mm.Takt_Pieces_x_1M,

			&mm.Job_Id, &mm.Product_Id, &mm.str_time_stamp, &mm.Pct_Slow_Cycle,
			&mm.Pct_Small_Stop, &mm.Pct_Full_Stop, &mm.Job_Desc, &mm.Ref_No,
			&mm.Ideal_Cycle_Hours_x_1M, &mm.Ideal_Cycle_Minutes_x_1M,
			&mm.Ideal_Cycle_Seconds_x_1M,
			&mm.Takt_Hours_x_1M, &mm.Takt_Minutes_x_1M,
			&mm.Takt_Seconds_x_1M,
			&mm.Scale_Total_Count_x_1K, &mm.Scale_Reject_Count_x_1K,
			&mm.Goal_Qty, &mm.Slow_Cycle__Treshold_Sequence_Display,
			&mm.Down_Rate_per_Hour, &mm.Run_Rate_per_Hour,
			&mm.Switch_from_Run_to_Down_based_on, &mm.Switch_from_Down_to_Run_based_on, 
			&mm.Switch_from_Down_to_Run_Count, //&mm.Switch_from_Down_to_Run_counter_of, 
			&mm.Rate_per_Hour_counter_of, &mm.Remote_Code, &mm.Is_Deleted); err != nil {
			return "", array_Id, array_job, nil, err
		}
		if mm.Switch_from_Run_to_Down_based_on < 0 || mm.Switch_from_Run_to_Down_based_on > 1 {
			mm.Switch_from_Run_to_Down_based_on = 0
		}
		if mm.Switch_from_Down_to_Run_based_on < 0 || mm.Switch_from_Down_to_Run_based_on > 2 {
			mm.Switch_from_Down_to_Run_based_on = 0
		}
		if mm.Switch_from_Down_to_Run_Count < 0 || mm.Switch_from_Down_to_Run_Count > 100 {
			mm.Switch_from_Down_to_Run_Count = 10
		}
		if mm.Rate_per_Hour_counter_of < 0 || mm.Rate_per_Hour_counter_of > 1 {
			mm.Rate_per_Hour_counter_of = 0
		}
		if mm.Slow_Cycle__Treshold_Sequence_Display > 100 { mm.Slow_Cycle__Treshold_Sequence_Display = 100 }

		array_job = append(array_job, mm)
		fmt.Println(mm)
	}
	return qryD[1], array_Id, array_job, nil, nil
}