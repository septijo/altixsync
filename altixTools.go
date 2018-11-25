package main

import (
	"bufio"
	"crypto/sha256"
	//"database/sql"
	"encoding/hex"
	//"errors"

	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

////////////
func string_tenary(cond bool, true_return string, false_return string) (string) {
	////////////
		if cond { return true_return } else { return false_return }
	}

	
func np__ganti_lineFeed_jadi_spasi__trim_semua_space__padR_dgn_0(input string, pad_R_len uint) []byte {
	//fmt.Println("#0",[]byte(input))
	input = strings.Replace(input, "\r\n", " ", -1)
	input = strings.Replace(input, "\n", " ", -1)
	input = strings.TrimSpace(input)
	ret := make([]byte, pad_R_len)
	//fmt.Println("#1",[]byte(input))
	var i uint
	for i = 0; i < uint(len(input)); i++ {
		ret[i] = input[i]
	}
	//fmt.Println("#2",ret)
	for ; i < pad_R_len-1; i++ {
		ret[i] = 0
	}
	//fmt.Println("#3",ret)
	return ret
}

//////////////////////////////////////////////////////////////////
//                                                              //
func np_get_yymdhns_from_byte_slice(msg []byte) string {
	//                                                              //
	//////////////////////////////////////////////////////////////////
	ret := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", int(msg[0])*100+int(msg[1]),
		msg[2], msg[3]&31, msg[4], msg[5], msg[6])
	if msg[3]&224 > 0 {
		dow := [8]string{"", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		ret += " " + dow[(msg[3]>>5)&7]
	}
	return ret
}

//////////////////////////////////////////////////////////////////
//                                                              //
func np_get_array_yymdhns_from_byte_slice(msg []byte) []string {
	//                                                              //
	//////////////////////////////////////////////////////////////////
	ret := make ([]string,2)
	ret[0] = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", int(msg[0])*100+int(msg[1]),
		msg[2], msg[3]&31, msg[4], msg[5], msg[6])
	if msg[3]&224 > 0 {
		dow := [8]string{"", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		ret[0] += " " + dow[(msg[3]>>5)&7]
	}
	ret[1] = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", int(msg[7])*100+int(msg[8]),
		msg[9], msg[10]&31, msg[11], msg[12], msg[13])
	return ret
}

//////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                          //
func np_write_datetime120_to_byte_slice7(datetime120 string, e *[]byte, startPos int) {
	//                                                                                          //
	//////////////////////////////////////////////////////////////////////////////////////////////
	thn, _ := strconv.ParseInt(datetime120[0:4], 10, 0)
	bln, _ := strconv.ParseInt(datetime120[5:7], 10, 0)
	tgl, _ := strconv.ParseInt(datetime120[8:10], 10, 0)
	jam, _ := strconv.ParseInt(datetime120[11:13], 10, 0)
	mnt, _ := strconv.ParseInt(datetime120[14:16], 10, 0)
	dtk, _ := strconv.ParseInt(datetime120[17:19], 10, 0)
	(*e)[startPos+0] = byte(thn / 100)
	(*e)[startPos+1] = byte(thn % 100)
	(*e)[startPos+2] = byte(bln)
	(*e)[startPos+3] = byte(tgl)
	(*e)[startPos+4] = byte(jam)
	(*e)[startPos+5] = byte(mnt)
	(*e)[startPos+6] = byte(dtk)
}

///////////////////////////////////////////////////////////////////////////////////////////
//                                                                                       //
func nc_write_uint64_to_byte_slice__LSB_to_MSB(d uint64, e *[]byte, startPos int) {
	//                                                                                       //
	///////////////////////////////////////////////////////////////////////////////////////////
	(*e)[startPos+0] = byte(d % 256)
	(*e)[startPos+1] = byte(d << 48 >> 56)
	(*e)[startPos+2] = byte(d << 40 >> 56)
	(*e)[startPos+3] = byte(d << 32 >> 56)
	(*e)[startPos+4] = byte(d << 24 >> 56)
	(*e)[startPos+5] = byte(d << 16 >> 56)
	(*e)[startPos+6] = byte(d << 8 >> 56)
	(*e)[startPos+7] = byte(d << 0 >> 56)
}

///////////////////////////////////////////////////////////////////////////////////////////
//                                                                                       //
func nc_write_uint32_to_byte_slice__LSB_to_MSB(d uint32, e *[]byte, startPos int) {
	//                                                                                       //
	///////////////////////////////////////////////////////////////////////////////////////////
	(*e)[startPos+0] = byte(d % 256)
	(*e)[startPos+1] = byte(d << 16 >> 24)
	(*e)[startPos+2] = byte(d << 8 >> 24)
	(*e)[startPos+3] = byte(d << 0 >> 24)
}

///////////////////////////////////////////////////////////////////////////////////////////
//                                                                                       //
func nc_write_uint16_to_byte_slice__LSB_to_MSB(d uint16, e *[]byte, startPos int) {
	//                                                                                       //
	///////////////////////////////////////////////////////////////////////////////////////////
	(*e)[startPos+0] = byte(d % 256)
	(*e)[startPos+1] = byte(d >> 8)
}

/////////////////////////////////////////////
//                                         //
func npSha256(data string) string {
	//                                         //
	/////////////////////////////////////////////
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

/////////////////////////////////////////////
//                                         //
func npSha256_from_byte_arr(data []byte) string {
	//                                         //
	/////////////////////////////////////////////
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

/////////////////////////////////////////////
//                                         //
func npSha256_from_int_arr(data []int) string {
	//                                         //
	/////////////////////////////////////////////
	h := sha256.New()
	var b = make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		b[i] = byte(data[i])
	}
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

//////////
func read_n_bytes(conn net.Conn, jByte int) ([]byte, int, error) {
	var readCount = 0
	result := make([]byte, jByte)
	//var err = 0

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {

		byte_read, err := bufio.NewReader(conn).ReadByte()
		if err != nil {
			//break
			return result, readCount, err
		}
		result[readCount] = byte_read
		readCount++
		if readCount == jByte {
			return result, readCount, nil
		}
	}
}
