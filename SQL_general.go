package main

import (
	//"time"
	// utk bisa baca file .ini harus "go get github.com/c4pt0r/cfg"
	"github.com/c4pt0r/cfg"

	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
	"strings"
)

var db *sql.DB

///////////////////////////
//                       //
func MSSQL_open() {
	//                       //
	///////////////////////////
	var conf = cfg.NewCfg("./altixSync.ini")
	conf.Load()
	var (
		sqlIpAddress, _ = conf.ReadString("ipAddress", "")
		sqlPort, _      = conf.ReadString("port", "")
		sqlUserId, _    = conf.ReadString("userId", "")
		sqlPassword, _  = conf.ReadString("password", "")
	)
	fmt.Println("server=" + sqlIpAddress + ";port=" + sqlPort + ";user id=" + sqlUserId + ";password=" + sqlPassword)
	var err error
	db, err = sql.Open("mssql", "server="+sqlIpAddress+";port="+sqlPort+
		";user id="+sqlUserId+";password="+sqlPassword+";database=Altix")
	if err != nil {
		fmt.Println("func open_MSSQL() : " + err.Error())
	}
}

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
func MSSQL_get__tenantId__deviceId__timeStampS__of_mcuId(mcu_Id string,
	storage_Id string) (tenant_Id int64, altix_device_Id int64,
	machine_Id int64, aTime [][]string, err error) {
	//                                                                         //
	/////////////////////////////////////////////////////////////////////////////
	//var aTime = make([]time.Time, 2)
	//for i := 0; i < 2; i++ {
	//	aTime[i] = time.Time{}
	//}
	//var aTime = make([]string, 7)
	aTime = make([][]string, 11)
	for i := range aTime {
		aTime[i] = make([]string, 2)
	}
	//for i := 0; i < 2; i++ { aTime[i] = "" }

	sqlQuery := "select isnull(MP.Tenant_Id,-1), AD.Altix_Device_Id, " +
		"  isnull(MM.Machine_Id,-1), convert(varchar(19), MM.Time_Stamp, 120) " +
		//"from mstMachine as MM " +
		//"left join altix_Devices as AD on MM.altix_Device_Id = AD.altix_Device_Id " +
		"from altix_Devices as AD " +
		"left join mstMachine as MM on MM.Altix_Device_Id = AD.Altix_Device_Id " +
		"left join mstPlant_Line as MPL on MPL.Plant_Line_Id = MM.Plant_Line_Id " +
		"left join mstPlant_Area as MPA on MPA.Plant_Area_Id = MPL.Plant_Area_Id " +
		"left join mstPlant as MP on MPA.Plant_Id = MP.Plant_Id " +
		"where AD.Processor_Id = '" + mcu_Id + "'"
	rows, err := db.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error() + "\n" + sqlQuery)
		return -1, -1, -1, aTime, err
	}

	if rows.Next() {
	} else {
		//fmt.Println("Altix Device with MCU ID [" + mcu_Id + "] NOT FOUND !!")
		sqlCmd := "insert into altix_devices (Processor_Id, Storage_Id, Altix_Model, " +
			"Description, PCB_Version, Processor_Type, Storage_Type, Display_Height_Pixel, " +
			"Display_Width_Pixel, Append_DateTime, Time_Stamp) select '" + mcu_Id + "', '" +
			storage_Id + "', '', 'insert by ALtixSync','', '', '', -1, -1, getDate(), getDate()"
		_, err := db.Exec(sqlCmd)
		if err != nil {
			fmt.Println(err.Error() + "\n" + sqlCmd)
		}
	}

	var device_setting_timeStamp sql.NullString
	err = rows.Scan(&tenant_Id, &altix_device_Id, &machine_Id, &device_setting_timeStamp)
	if err != nil {
		fmt.Println("#get__ti_di_ts__of_mcuid A : rows.Scan " + err.Error())
		return -1, -1, -1, aTime, err
	}
	aTime[0][0] = device_setting_timeStamp.String
	fmt.Println("\n     SQL Tenant_Id =", tenant_Id, "\n     SQL Device setting                  = "+aTime[0][0])

	if tenant_Id == -1 {
		return -1, -1, -1, aTime, errors.New("#get__ti_di_ts__of_mcuid : Altix Device [" +
			mcu_Id + "] not linked to any Machine yet.")
		//return -1, -1, aTime, err
	}

	// sekarang baca semua setting dan master utk tau Time Stamp yg paling baru
	aTime[1], _ = get_LastTimeStamp("mstMetric              ", "       ", tenant_Id)
	aTime[2], _ = get_LastTimeStamp("mstReason              ", "Down   ", tenant_Id)
	aTime[3], _ = get_LastTimeStamp("mstReason              ", "Standby", tenant_Id)
	aTime[4], _ = get_LastTimeStamp("mstReason              ", "Setup  ", tenant_Id)
	aTime[5], _ = get_LastTimeStamp("mstShift               ", "       ", tenant_Id)
	aTime[6], _ = get_LastTimeStamp("mstBreak               ", "       ", tenant_Id)
	aTime[7], _ = get_LastTimeStamp("TimeSchedule_Shift     ", "       ", machine_Id)
	aTime[8], _ = get_LastTimeStamp("TimeSchedule_ShiftBreak", "       ", machine_Id)
	aTime[9], _ = get_LastTimeStamp("mstProduct             ", "       ", tenant_Id)
	aTime[10], _ = get_LastTimeStamp("job                    ", "       ", machine_Id)

	return tenant_Id, altix_device_Id, machine_Id, aTime, nil
}

////////////
func get_LastTimeStamp_Query(ndx int, jenis string, reasonType string, 
	tenant_or_machine_Id int64) (qry string) {
///////////

	field_Id := "Tenant_Id"
	if strings.TrimSpace(jenis) == "TimeSchedule_ShiftBreak" {
		qry = "select top 1 convert(varchar(19), TS_SB.Time_Stamp, 120) " +
			"from TimeSchedule_ShiftBreak as TS_SB " +
			"left join TimeSchedule_Shift as TS_S " +
			"  on TS_SB.TimeSchedule_Shift_Id = TS_S.TimeSchedule_Shift_Id " +
			"where TS_SB.Is_Deleted = " + strconv.Itoa(ndx) + " and TS_S.Machine_Id = " + strconv.FormatInt(tenant_or_machine_Id, 10) +
			" order by convert(varchar(19), TS_SB.Time_Stamp, 120) desc"
		field_Id = "TS_S.Machine_Id" // khusus "if" yg ini, utk  DEBUG_MSG doang
	} else {
		if strings.TrimSpace(jenis) == "TimeSchedule_Shift" ||
			strings.TrimSpace(jenis) == "job" ||
			strings.TrimSpace(jenis) == "mstMachine" {
			field_Id = "machine_Id"
		}

		var cekDelete string
		cekDelete = "Is_Deleted = " + strconv.Itoa(ndx) + " and "
		if strings.TrimSpace(jenis) == "mstMetric" {
			cekDelete = ""
		}
		qry = "select top 1 convert(varchar(19), Time_Stamp, 120) " +
			"from " + jenis + " where " + cekDelete + field_Id + " = " + strconv.FormatInt(tenant_or_machine_Id, 10)
		if strings.TrimSpace(reasonType) != "" {
			qry += " and Reason_Type = '" + strings.TrimSpace(reasonType) + "' "
		}
		qry += " order by convert(varchar(19), Time_Stamp, 120) desc"
	}
	return qry
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
//
func get_LastTimeStamp(jenis string, reasonType string, tenant_or_machine_Id int64) ([]string, error) {
	//
	////////////////////////////////////////////////////////////////////////////////////////////////////////
	//field_Id := "Tenant_Id"
	var qry = make([]string, 2)
	var result = make([]string, 2)
	for i := 0; i < 2; i++ {
		qry[i] = get_LastTimeStamp_Query(i,jenis, reasonType, tenant_or_machine_Id)
		if qry[i] == "" { continue }
		rows, err := db.Query(qry[i])
		if err != nil {
			fmt.Println("func get_LastTimeStamp(" + strings.TrimSpace(jenis) + ") error: " + err.Error() + "\n" + qry[i])
			return result, err
		}
		if rows.Next() {
		} else {
			//fmt.Println("func get_TimeStamp (Is_Deleted = " + strconv.Itoa(i) + ")" + jenis + " for " + field_Id + " [" +
			//	strconv.FormatInt(tenant_or_machine_Id, 10) + "] NOT FOUND !!")
			result[i] = "1900-01-01 00:00:00"
			continue
			//return "", errors.New("#get__timeStamp : 0 rows.")
		}
		var ret sql.NullString
		err = rows.Scan(&ret)
		if err != nil {
			fmt.Println("func get_TimeStamp : rows.Scan " + err.Error())
			return result, err
		}
		//fmt.Println("     SQL " + jenis + " " + reasonType + " = " + ret.String)
		result[i] = ret.String
	}
	fmt.Println("     SQL " + jenis + " " + reasonType + " =", result)
	return result, nil
}

