package main

import (
	"bufio"
	"math/rand"

	//"flag" // utk urusan ms-sql
	"fmt"
	//"io"
	//"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// utk bisa baca file json sebagai config harus "go get github.com/tkanos/gonfig"
	//"github.com/tkanos/gonfig"

	// utk bisa baca file .ini harus "go get github.com/c4pt0r/cfg"
	// "github.com/c4pt0r/cfg"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	//CONN_HOST = "localhost"
	CONN_HOST  = ""
	CONN1_PORT = "7100"
	CONN2_PORT = "7000"
	CONN_TYPE  = "tcp"
)

const MAX int = 1000000

type altixDevice struct {
	conn                                   net.Conn
	sliceUsed                              bool
	slice_lock                             sync.Mutex
	tenant_Id, altix_device_Id, machine_Id int64
	mcu_Id                                 string

	socket_perlu_diclose bool

	// altix_device setting/config Time Stamp
	dvc_device_setting_timeStamp          string
	dvc_master_metric_timeStamp           string
	dvc_master_downReason_timeStamp       []string
	dvc_master_standbyReason_timeStamp    []string
	dvc_master_setupReason_timeStamp      []string
	dvc_master_shift_timeStamp            []string
	dvc_master_break_timeStamp            []string
	dvc_timeSchedule_Shift_timeStamp      []string
	dvc_timeSchedule_ShiftBreak_timeStamp []string
	dvc_master_product_timeStamp          []string
	dvc_job_timeStamp                     []string

	// MS-SQL Time Stanp :
	SQL_device_setting_timeStamp          string //time.Time
	SQL_master_metric_timeStamp           string
	SQL_master_downReason_timeStamp       []string
	SQL_master_standbyReason_timeStamp    []string
	SQL_master_setupReason_timeStamp      []string
	SQL_master_shift_timeStamp            []string
	SQL_master_break_timeStamp            []string
	SQL_timeSchedule_Shift_timeStamp      []string
	SQL_timeSchedule_ShiftBreak_timeStamp []string
	SQL_master_product_timeStamp          []string
	SQL_job_timeStamp                     []string
}

var altixDevices = make([]altixDevice, MAX) // edan gan ..1 juta client !!

var global_lock sync.Mutex

//////////////////////////////////////////////////////////
//                                                      //
func findSlice_by_tenantId(tenant_id int64) int {
	//                                                      //
	//////////////////////////////////////////////////////////
	global_lock.Lock()
	for i := 0; i < MAX; i++ {
		if altixDevices[i].sliceUsed && altixDevices[i].tenant_Id == tenant_id {
			global_lock.Unlock()
			return i
		}
	}
	global_lock.Unlock()
	return -1
}

////////////////////////////////////////
//                                    //
func cariSliceYangKosong__dan__kuasai(conn net.Conn) int {
	//                                    //
	////////////////////////////////////////
	global_lock.Lock() // // karena slice ini akan di kuasai, pakai global_lock
	for i := 0; i < MAX; i++ {
		if !altixDevices[i].sliceUsed {
			altixDevices[i].conn = conn
			altixDevices[i].sliceUsed = true
			altixDevices[i].socket_perlu_diclose = false
			global_lock.Unlock()
			return i
		}
	}
	global_lock.Unlock()
	return -1
}

const ymdhnsDateTimeFmt = "2006-01-02 15:04:05"

var recvLengthErrorCount = 0

/////////////////////
//                 //
func main() {
	//                 //
	/////////////////////
	rand.Seed(time.Now().UTC().UnixNano())
	MSSQL_open() // ada di sql_general.go
	defer db.Close()

	go listen_for_altix_dotNet() // ada dibawah

	// Listen for incoming altixSync connections
	l1, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN2_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l1.Close()
	fmt.Println("Listening for altixDevice on *:" + CONN2_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l1.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// cari yg kosong
		var pos = cariSliceYangKosong__dan__kuasai(conn)
		if pos >= 0 {
			// Handle connections in a new goroutine.
			go handleRequest_altixDevice_sync(conn, pos) // ada di altixProtocol.go
		} else {
			fmt.Println("Server is Full (all altixDevice slice was taken), closing socket.")
			conn.Close()
		}
	}
}

//////////////////////////////////////
//                                  //
func closeConn(slicePos int) {
	//                                  //
	//////////////////////////////////////
	global_lock.Lock() // karena slice ini akan di reset, pakai global_lock
	altixDevices[slicePos].sliceUsed = false
	altixDevices[slicePos].socket_perlu_diclose = false
	altixDevices[slicePos].conn.Close()
	altixDevices[slicePos].conn = nil
	global_lock.Unlock()
}

func listen_for_altix_dotNet() {
	// Listen for incoming altixDotNet connections
	l2, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN1_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l2.Close()
	fmt.Println("Listening for altixDotNet on *:" + CONN1_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l2.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest_altixDotNet_notif(conn)
	}
}

/////////////////////////////////////////////////////////////
//                                                         //
func handleRequest_altixDotNet_notif(conn net.Conn) {
	//                                                         //
	/////////////////////////////////////////////////////////////
	fmt.Println("\n"+time.Now().Format(ymdhnsDateTimeFmt), "NOTIF", conn.RemoteAddr(), "connected.")
	message, _ := bufio.NewReader(conn).ReadString('\n')

	if len(message) < 20 {
		conn.Close()
		return
	}

	fmt.Println("rcvd message:", message)

	machineId, _ := strconv.ParseInt(message[5:15], 10, 32)
	// *GET /0123456789/DeviceSetting HTTP/1.1 *
	// *GET /0123456789/blabla HTTP/1.1 *

	ndx := strings.Index(message, " HTTP")
	//if ndx == -1 { ndx = len(message) - 1 }
	if ndx == -1 {
		conn.Close()
		return
	}
	data := message[16:ndx]

	fmt.Println("Rcv", len(message), "byte", conn.RemoteAddr(), "*", string(message),
		"* machineId /data / ndx=", machineId, "/", data, len(message), "*")
	conn.Write([]byte("X 201 OK\n"))
	conn.Close()
	//var send = ",NU:" + data + ","
	// notif soket altixDevice ybs sesuai altixId masing2
	global_lock.Lock()
	for i := 0; i < MAX; i++ {
		if altixDevices[i].sliceUsed && altixDevices[i].machine_Id == machineId {
			time.Sleep(2 * time.Second)
			if strings.EqualFold(data, "DeviceSetting") {
				SQL_device_setting_timeStamp, _ := get_LastTimeStamp("mstMachine", "", machineId)
				altixDevices[i].SQL_device_setting_timeStamp = SQL_device_setting_timeStamp[0]
			}

			//if strings.EqualFold(data, "DEL_TimeSchedule_Shift") {
			//	altixDevices[i].timeSchedule_Shift_timeStamp = "1900-01-01 00:00:00"
			//}
			if strings.EqualFold(data, "TimeSchedule_Shift") {
				altixDevices[i].SQL_timeSchedule_Shift_timeStamp, _ = get_LastTimeStamp("TimeSchedule_Shift", "", machineId)
			}

			//if strings.EqualFold(data, "DEL_TimeSchedule_ShiftBreak") {
			//	altixDevices[i].timeSchedule_ShiftBreak_timeStamp = "1900-01-01 00:00:00"
			//}
			if strings.EqualFold(data, "TimeSchedule_ShiftBreak") {
				altixDevices[i].SQL_timeSchedule_ShiftBreak_timeStamp, _ = get_LastTimeStamp("TimeSchedule_ShiftBreak", "", machineId)
			}
		}
	}
	global_lock.Unlock()
}
