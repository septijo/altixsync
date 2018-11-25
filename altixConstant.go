package main

// data yang di send adalah :
//    byte #0 = 0x00 (protocol version)
//    byte #1 = randomGarbage[0]
//    byte #2 = randomGarbage[1]
//    byte #3 = randomGarbage[2]
//    byte #4 = randomGarbage[3]
//    byte #5 s/d (#4+garbageCount) = randomGarbage sebanyak garbageCount
//       garbageCount di hitung dari jumlah elemen #0 dari semua altix_random_token[randomGarbage[0..3]]

//    mulai byte (#5+garbageCount), yaitu 20 byte selanjutnya / terakhir adalah cpuId_storageId
//    byte #12,#13 dari cpuId_storageId adalah sum16 dari random garbage
//    byte #14,#15 dari cpuId_storageId adalah sum16 dari cpuId_storageId
//    contoh cpuId_Stroage_id
//        arr_cpuId_storageId := [20]byte{11,12,13,14,15,16,17,18,19,20,21,22,0,0,0,0,57,58,59,60}
//        sum dari cpuId_storageId = 432, 432 basis 10 = 0x01B0 basis 16, jadi byte #14 = 0x01, byte #15 = 0xB0 (176)

type altix_errorMessage struct {
	error_Id uint32
	error_Desc string
}

var altix_errorMessages = []altix_errorMessage { { 
	0x00000700, "_INCOMING_DEVICE_SETTING__WRONG_RECV_SIZE_"},	{ 
	0x00010700, "_INCOMING_DEVICE_SETTING__MACHINE_ID_IS_ZERO_"},	{ 
	0x00020700, "_INCOMING_DEVICE_SETTING__INVALID_SUFFIX_"},	{ 
	0x00050700, "_INCOMING_DEVICE_SETTING__INVALID_YYMDHNS_TIME_STAMP_"},	{ 
	0x00060700, "_INCOMING_DEVICE_SETTING__INVALID_YYMDHNS_NEXT_TIMESCHED_"},	{ 
	0x00070700, "_INCOMING_DEVICE_SETTING__INVALID_YYMDHNS_OCCASION_START_"},	{ 
	0x00080700, "_INCOMING_DEVICE_SETTING__INVALID_YYMDHNS_OCCASION_END_"},	{ 
	//0x11090700, "_INCOMING_DEVICE_SETTING__METRIC_ID_S1_R1_IS_INVALID_"},	{ 
	//0x12090700, "_INCOMING_DEVICE_SETTING__METRIC_ID_S1_R2_IS_INVALID_"}, {

	0x00000b00, "_INCOMING_MST_DOWN_REASON__WRONG_RECV_SIZE_"} , {
	0x00010b00, "_INCOMING_MST_DOWN_REASON__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00000c00, "_INCOMING_MST_STANDBY_REASON__WRONG_RECV_SIZE_"} , {
	0x00010c00, "_INCOMING_MST_STANDBY_REASON__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00000d00, "_INCOMING_MST_SETUP_REASON__WRONG_RECV_SIZE_"} , {
	0x00010d00, "_INCOMING_MST_SETUP_REASON__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00010e00, "_INCOMING_MST_SHIFT__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00010f00, "_INCOMING_MST_BREAK__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00011000, "_INCOMING_TIME_SCHEDULE_SHIFT__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00011100, "INCOMING_TIME_SCHEDULE_SHIFT_BREAK__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00011200, "_INCOMING_MST_PRODUCT__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0x00001300, "_INCOMING_JOB__WRONG_RECV_SIZE_"} , {
	0x00011300, "_INCOMING_JOB__INVALID_YYMDHNS_LAST_DELETE_"} , {

	0, ""} , {
	0, ""} , {
	0xFFFFFFFF, "dummy buat nutup agar gampang copy paste"} }


var altix_random_token = [256][28]int{
	{28,23,7,14,3,24,2,21,19,27,15,18,12,5,20,26,1,8,17,22,25,6,9,10,11,4,16,13}, //0
	{10,9,27,26,17,14,24,20,12,18,25,13,6,8,1,2,28,21,5,11,22,4,7,3,19,23,16,15},
	{23,12,27,19,24,15,4,9,5,25,10,17,26,2,18,16,8,13,21,20,3,28,11,6,1,14,22,7},
	{18,6,20,19,14,9,24,4,26,10,8,11,1,2,16,21,17,3,23,22,12,7,13,15,5,28,25,27},
	{17,4,23,1,24,12,16,2,21,19,27,14,7,26,20,9,25,13,22,11,18,8,5,15,10,3,28,6},
	{3,4,7,13,22,27,24,12,20,16,5,15,21,14,8,17,26,19,9,25,23,28,11,6,18,10,2,1}, //5
	{13,23,24,21,1,2,14,25,6,27,5,10,3,17,8,11,7,9,12,4,22,15,16,26,19,28,20,18},
	{26,14,24,18,20,1,15,28,8,21,17,11,23,25,3,13,16,27,5,6,12,19,10,2,4,22,9,7},
	{16,6,26,25,9,1,14,3,12,4,8,19,5,2,21,7,20,28,23,18,17,10,13,15,11,27,22,24},
	{17,9,1,7,25,14,26,21,2,28,4,3,6,18,16,19,12,13,10,23,24,15,11,27,8,5,22,20},
	{23,20,26,3,12,2,22,27,9,11,14,17,1,5,19,8,6,18,21,15,25,10,24,28,16,13,4,7}, //10
	{16,24,21,23,2,15,7,25,27,13,12,4,28,19,1,5,22,17,6,3,10,20,8,9,11,14,26,18},
	{21,3,13,26,27,20,28,7,1,15,4,23,6,19,8,24,25,16,2,18,14,9,12,17,22,5,11,10},
	{14,4,2,8,7,27,19,22,28,25,5,20,26,15,21,6,9,17,23,18,3,24,10,16,1,12,11,13},
	{25,1,9,23,5,28,24,22,11,13,4,3,15,2,16,17,14,6,27,20,8,26,10,7,19,12,21,18},
	{23,25,14,19,15,8,18,16,13,20,17,5,24,9,3,4,27,7,12,6,10,28,2,1,22,11,21,26}, //15
	{16,5,19,23,14,2,1,13,24,20,7,3,12,11,28,22,15,6,26,18,9,25,17,8,27,21,10,4},
	{17,20,11,6,5,15,28,22,9,21,16,10,25,26,1,3,24,27,12,7,2,23,13,18,19,4,8,14},
	{25,1,19,16,2,21,10,11,7,22,4,24,15,27,3,18,20,6,17,12,13,14,8,5,9,23,28,26},
	{23,9,21,3,15,6,5,25,7,14,16,26,27,28,10,11,17,24,22,8,2,19,13,18,4,1,20,12},
	{16,25,10,3,21,1,7,26,13,17,6,5,11,27,8,18,28,14,22,12,23,24,9,15,19,20,2,4}, //20
	{23,11,25,22,8,17,10,5,26,20,14,18,6,13,19,21,27,4,12,2,28,3,16,24,9,7,1,15},
	{16,19,12,26,24,18,1,2,3,9,8,13,7,21,6,4,14,25,23,10,20,15,22,11,17,27,28,5},
	{26,21,2,5,17,16,14,9,4,20,15,11,8,27,6,10,22,13,7,19,3,12,24,1,23,25,18,28},
	{12,21,18,7,15,1,26,11,6,20,4,8,13,19,9,10,16,5,23,24,17,3,14,22,27,25,28,2},
	{4,26,28,1,15,7,25,3,14,19,8,21,27,6,12,16,9,23,5,24,2,22,11,10,17,13,20,18}, //25
	{14,18,5,10,19,22,26,27,2,13,1,28,6,12,25,23,17,3,16,21,8,11,7,15,20,24,4,9},
	{1,18,2,7,6,17,21,27,14,15,19,8,3,16,9,20,25,23,28,13,5,22,12,10,4,24,11,26},
	{1,9,26,12,14,24,11,13,10,3,23,6,18,21,2,5,8,4,19,16,15,28,17,7,22,27,20,25},
	{17,24,5,23,4,20,27,19,2,13,7,16,14,18,22,1,15,25,28,3,8,10,12,26,21,6,11,9},
	{19,20,4,5,17,13,22,18,6,1,23,21,16,26,14,12,25,7,8,24,27,15,28,2,3,9,10,11}, //30
	{22,21,8,16,11,26,15,3,9,6,7,13,2,27,5,1,12,20,19,28,24,17,4,18,10,25,23,14},
	{4,28,1,15,13,27,26,21,5,11,24,14,2,23,25,7,17,12,6,16,20,18,8,10,19,22,9,3},
	{10,8,23,27,5,20,4,22,13,14,24,7,17,12,19,15,3,1,6,11,25,2,9,18,26,16,28,21},
	{16,27,26,12,23,11,28,5,20,21,2,15,13,10,25,18,6,9,17,22,8,3,4,14,1,7,19,24},
	{23,24,9,6,16,17,25,5,13,7,3,26,4,8,27,14,22,2,20,15,28,10,11,18,12,1,19,21}, //35
	{7,12,24,8,27,19,28,26,21,17,15,25,6,14,10,22,3,16,5,1,13,23,11,18,20,9,4,2},
	{27,22,26,28,16,2,1,25,15,11,20,19,17,13,6,5,7,21,3,12,24,8,23,4,10,14,9,18},
	{18,12,3,9,21,26,23,4,10,1,5,28,24,27,17,16,11,14,2,19,25,22,6,7,8,13,20,15},
	{15,18,17,9,13,8,24,25,2,1,12,19,22,4,5,3,11,26,28,6,23,20,10,14,16,7,21,27},
	{13,3,4,17,23,24,5,11,26,28,18,21,2,20,19,8,14,16,12,10,22,9,25,1,27,7,6,15}, //40
	{24,2,27,5,1,22,18,23,9,4,21,12,20,7,14,15,6,3,19,28,11,8,25,16,10,13,26,17},
	{22,9,28,26,3,14,23,17,4,15,8,16,25,5,7,21,19,27,2,1,12,24,6,18,13,11,10,20},
	{22,13,10,12,23,19,2,14,21,28,26,11,20,5,17,8,1,27,7,24,3,6,15,4,18,9,16,25},
	{13,12,2,20,24,19,10,15,18,23,6,17,16,8,7,22,3,4,11,9,5,28,14,25,27,26,21,1},
	{18,2,19,16,11,13,20,10,9,14,21,7,15,12,25,17,23,5,6,28,24,4,8,3,1,26,22,27}, //45
	{8,5,1,28,16,11,23,27,19,17,18,25,14,9,24,26,13,7,20,4,3,12,15,22,10,6,21,2},
	{14,8,12,26,28,13,22,9,6,3,7,17,10,16,18,20,23,15,27,4,25,11,5,21,1,19,2,24},
	{3,9,27,22,1,13,19,11,20,8,5,7,10,18,4,16,21,15,14,24,12,6,28,23,25,2,17,26},
	{24,20,19,16,13,25,1,8,18,23,2,6,14,5,27,10,9,15,17,4,26,22,7,3,12,28,21,11},
	{10,24,16,1,27,20,6,26,2,11,17,19,22,8,18,21,15,7,12,3,23,14,9,13,5,25,28,4}, //50
	{20,8,3,12,6,14,10,5,17,24,26,25,11,2,16,4,1,9,13,15,18,28,23,27,19,21,22,7},
	{10,12,16,24,28,23,5,14,20,8,13,25,3,11,17,7,18,21,15,2,26,19,9,27,1,6,22,4},
	{25,13,2,20,1,11,14,21,26,19,16,22,9,24,23,28,17,6,7,3,4,8,10,5,27,12,15,18},
	{20,8,11,6,9,24,16,17,18,4,25,26,7,10,13,12,28,21,22,15,3,14,1,19,2,5,23,27},
	{25,24,19,1,20,6,4,16,14,28,13,9,8,18,21,12,23,11,2,7,26,27,10,3,15,5,17,22}, //55
	{8,12,1,11,23,5,15,4,27,13,22,19,14,17,10,2,9,6,25,26,21,24,18,16,3,28,20,7},
	{16,7,15,17,9,26,4,8,22,25,13,19,1,2,11,27,5,14,3,28,23,10,12,18,21,6,20,24},
	{5,22,27,1,6,10,13,15,11,3,16,26,28,4,2,19,9,24,25,23,12,8,14,20,17,7,21,18},
	{8,2,5,15,13,22,27,16,3,21,19,6,26,20,18,14,23,28,4,24,25,10,1,11,17,9,12,7},
	{6,12,2,3,18,28,15,23,10,20,1,21,9,19,16,7,17,22,14,4,13,5,8,11,27,26,24,25}, //60
	{5,25,17,19,4,20,9,28,10,11,21,8,14,27,2,24,15,18,22,3,7,12,13,1,6,26,23,16},
	{26,19,20,23,6,27,16,22,13,9,18,2,12,17,10,7,5,21,28,3,4,1,25,8,15,11,14,24},
	{26,23,22,19,2,28,1,7,3,15,25,10,24,9,14,18,17,4,16,11,27,20,12,6,8,5,13,21},
	{2,23,18,11,24,6,1,14,15,16,22,4,25,28,9,26,13,19,17,7,20,21,3,10,5,27,8,12},
	{26,6,22,17,2,24,27,28,7,5,9,12,18,13,11,23,21,15,1,16,10,20,4,19,3,14,25,8}, //65
	{19,14,27,26,18,16,4,21,15,24,13,12,1,7,20,10,9,17,5,2,11,6,3,28,8,25,22,23},
	{19,10,8,12,15,2,4,26,16,22,18,3,23,7,17,11,13,21,20,9,25,24,1,6,27,28,14,5},
	{4,13,22,19,15,6,26,9,11,7,25,1,5,14,12,3,24,23,17,18,8,21,27,28,20,2,16,10},
	{13,4,1,5,6,23,7,8,11,21,26,27,28,12,9,16,2,18,24,10,20,17,15,25,14,22,3,19},
	{14,1,20,11,3,21,6,23,12,15,7,27,17,18,19,28,8,25,2,16,13,4,10,26,22,24,9,5}, //70
	{20,15,6,21,11,8,14,27,25,9,19,22,26,3,10,5,7,12,17,24,1,2,16,23,4,13,28,18},
	{21,7,20,6,17,5,8,2,26,3,11,12,24,27,9,13,1,16,4,10,25,23,14,18,15,19,22,28},
	{9,11,8,4,15,18,20,3,6,10,22,14,5,16,19,17,26,28,13,27,7,12,23,2,25,21,24,1},
	{12,9,24,11,21,20,27,17,10,1,25,2,15,18,6,13,26,28,19,4,8,7,23,16,22,5,14,3},
	{14,24,13,5,9,26,27,10,17,16,25,28,7,6,2,4,11,21,23,18,20,3,12,19,1,8,22,15}, //75
	{23,1,8,5,17,12,15,28,13,6,27,21,4,7,9,22,14,16,3,24,10,26,18,19,20,25,2,11},
	{3,15,25,22,11,26,2,14,13,9,16,21,18,17,27,28,6,23,19,20,1,4,7,5,8,12,24,10},
	{11,19,16,13,23,27,12,7,15,20,10,5,8,1,22,25,28,14,2,3,4,21,9,18,24,17,26,6},
	{2,12,27,26,13,11,25,3,18,4,20,14,19,22,23,15,28,5,9,10,21,7,17,6,8,16,24,1},
	{1,12,9,23,2,25,27,24,19,8,10,5,14,13,6,7,16,3,17,20,11,22,26,28,15,18,4,21}, //80
	{2,22,6,15,4,26,24,1,7,3,17,21,18,11,19,5,28,20,8,14,27,23,9,12,10,13,16,25},
	{10,26,22,25,14,7,9,11,8,13,18,12,2,23,3,4,6,1,24,20,19,28,15,5,16,21,17,27},
	{5,17,23,10,15,8,2,27,6,26,14,19,9,25,12,16,18,4,20,13,24,28,3,7,1,11,21,22},
	{9,20,13,18,10,3,5,24,28,2,16,14,12,1,19,8,25,11,22,21,27,4,26,17,15,23,7,6},
	{6,8,24,25,4,23,5,1,13,7,16,10,11,12,3,26,19,9,18,28,15,22,17,20,21,27,2,14}, //85
	{12,4,8,23,17,7,18,26,16,21,6,24,11,9,22,25,10,15,27,3,5,19,1,2,20,14,13,28},
	{4,28,27,21,22,11,14,6,25,2,20,8,5,3,13,10,24,12,15,16,18,7,19,23,9,26,1,17},
	{9,17,26,27,23,7,15,25,2,28,1,6,13,5,21,12,11,18,4,20,3,19,8,24,16,14,22,10},
	{15,25,5,23,6,14,24,16,7,3,8,19,10,20,11,22,13,12,21,27,28,18,17,26,4,1,9,2},
	{16,21,15,28,9,23,14,1,20,3,10,4,5,8,6,7,24,22,17,25,11,26,12,18,27,19,2,13}, //90
	{16,9,18,26,2,5,8,14,13,15,11,19,25,6,10,12,24,23,21,20,4,28,1,7,17,27,3,22},
	{12,23,19,8,2,10,26,5,18,20,9,15,6,1,28,14,21,3,27,17,4,13,22,24,25,16,11,7},
	{23,18,3,5,11,7,15,28,26,13,12,4,20,9,8,16,19,1,6,25,21,27,10,17,24,2,22,14},
	{26,24,2,27,5,1,8,28,14,10,16,13,11,9,3,20,25,18,7,21,22,15,4,23,17,19,12,6},
	{25,15,1,6,11,26,16,21,4,24,27,20,17,28,12,7,9,13,5,23,18,19,10,8,3,14,22,2}, //95
	{15,23,7,3,14,19,6,21,1,16,8,4,20,24,18,28,5,27,2,17,10,11,22,26,12,13,9,25},
	{2,10,14,22,21,23,3,27,13,7,17,26,20,9,8,1,5,24,6,15,19,18,12,11,28,25,4,16},
	{7,10,8,13,28,27,14,20,23,11,25,15,6,16,4,19,24,9,5,17,22,2,26,21,12,18,1,3},
	{26,9,22,21,17,5,18,28,4,2,1,8,20,3,27,23,13,16,15,6,12,11,24,10,19,14,25,7},
	{23,17,1,6,14,2,13,18,27,11,12,15,8,7,9,28,24,25,20,22,10,26,3,4,19,21,16,5}, //100
	{23,17,1,6,14,2,13,18,27,11,12,15,8,7,9,28,24,25,20,22,10,26,3,4,19,21,16,5},
	{20,18,2,5,28,9,6,14,21,1,7,10,13,15,19,8,11,16,17,23,24,12,3,4,25,22,26,27},
	{5,16,7,21,12,20,17,9,11,14,26,18,3,27,23,19,4,10,8,13,22,2,24,1,28,6,25,15},
	{22,7,18,10,26,15,1,27,11,9,14,17,25,28,12,24,21,3,8,13,23,16,5,19,4,6,20,2},
	{2,4,10,8,26,16,3,19,17,24,6,13,12,27,5,28,20,25,18,23,9,1,22,14,7,11,21,15}, //105
	{15,22,7,26,12,16,25,6,1,14,8,27,9,4,21,23,13,20,17,2,24,10,28,11,19,5,18,3},
	{7,5,3,26,6,13,10,9,18,4,27,1,8,20,28,16,11,22,2,15,21,25,19,12,23,14,17,24},
	{24,18,5,26,19,27,22,15,8,7,25,17,6,12,10,11,14,20,28,4,3,9,1,16,21,13,23,2},
	{21,6,1,8,5,23,26,10,7,4,13,11,18,2,14,17,12,27,9,22,19,3,28,20,16,15,24,25},
	{24,23,16,11,21,26,1,28,20,4,17,25,13,19,3,18,5,12,6,15,10,7,8,9,14,2,22,27}, //110
	{10,13,17,26,7,12,15,14,21,25,2,22,4,18,8,16,11,1,3,28,9,24,20,23,19,5,27,6},
	{23,10,8,20,2,3,26,28,25,5,19,1,7,17,4,14,18,16,27,15,11,12,13,21,24,9,6,22},
	{7,28,14,24,20,22,3,23,18,9,16,4,6,8,11,27,1,26,2,17,13,12,19,21,15,10,25,5},
	{10,3,23,24,16,13,19,28,7,4,20,18,5,27,17,6,12,2,21,22,15,1,26,25,11,9,8,14},
	{27,23,15,19,3,6,14,4,10,9,5,1,26,11,17,18,28,13,16,20,22,24,8,7,12,25,21,2}, //115
	{16,10,17,8,14,19,18,7,12,15,5,20,21,22,6,28,9,2,4,23,25,13,26,11,27,3,1,24},
	{11,19,2,20,24,25,17,10,26,13,12,15,28,1,7,8,16,14,9,22,5,21,6,18,3,27,23,4},
	{10,5,4,13,24,17,6,11,14,12,9,1,16,19,28,18,25,7,2,3,26,22,20,15,21,8,27,23},
	{9,6,12,1,2,17,23,8,16,21,27,28,10,26,24,3,19,25,13,4,22,5,18,14,7,15,20,11},
	{10,23,21,2,28,17,13,14,16,22,5,18,11,9,7,26,15,19,20,8,24,1,4,25,27,6,3,12}, //120
	{26,23,4,2,20,11,13,19,6,28,24,10,1,3,15,14,9,22,27,16,7,5,12,18,21,17,25,8},
	{23,21,11,25,16,9,27,8,18,22,20,4,14,3,17,12,15,6,2,10,19,13,28,26,1,24,5,7},
	{4,21,25,9,12,8,22,27,6,18,7,24,15,1,5,3,23,14,17,11,26,10,2,16,28,13,20,19},
	{7,13,21,24,16,3,6,25,5,11,19,4,22,8,15,20,1,27,18,17,2,26,28,9,10,14,12,23},
	{22,13,14,11,12,24,23,7,20,25,4,26,6,15,3,8,27,9,18,16,2,10,28,17,19,5,1,21}, //125
	{8,7,15,18,20,25,11,27,21,19,13,16,10,2,9,22,26,6,24,28,17,1,12,14,23,5,3,4},
	{12,14,11,28,3,17,7,24,25,1,20,2,27,8,6,21,13,18,23,15,10,22,26,4,9,16,5,19},
	{17,2,4,13,8,6,7,9,24,11,10,22,15,25,18,5,1,14,16,3,21,28,20,19,27,26,23,12},
	{10,27,16,9,21,5,20,26,24,4,28,7,6,3,1,15,17,12,14,19,18,25,13,8,22,11,2,23},
	{6,7,17,12,2,19,16,28,1,23,18,25,15,24,20,4,14,22,10,27,8,11,5,13,9,3,21,26}, //130
	{1,26,18,12,23,15,8,7,2,24,17,10,16,6,11,4,21,22,3,9,5,13,14,20,28,19,27,25},
	{2,6,4,14,16,27,9,23,18,11,28,19,25,12,17,3,8,15,24,1,10,22,13,7,20,26,5,21},
	{26,17,22,24,4,18,15,7,2,9,1,14,16,25,3,23,21,28,13,8,11,5,6,27,20,12,19,10},
	{26,20,18,16,21,15,17,5,14,12,6,8,22,28,7,4,11,13,25,10,1,3,27,2,23,24,19,9},
	{2,25,1,7,22,28,27,13,19,4,3,17,15,12,9,8,5,6,20,10,16,18,14,24,11,21,23,26}, //135
	{9,18,7,3,15,6,10,5,8,27,24,21,11,1,12,14,20,4,13,2,23,19,28,22,26,25,16,17},
	{4,15,5,19,18,9,12,22,21,17,2,24,25,27,8,3,20,13,26,28,16,23,14,1,11,7,6,10},
	{21,14,1,7,18,25,28,16,3,24,20,19,4,5,11,27,22,6,26,9,2,10,15,23,8,17,12,13},
	{9,21,1,20,27,11,6,7,15,25,19,10,23,28,2,12,14,13,22,24,4,26,16,18,17,5,8,3},
	{28,18,12,15,21,22,27,24,4,23,9,13,6,3,19,17,2,16,10,25,14,8,5,7,26,20,1,11}, //140
	{18,17,24,22,6,1,23,3,12,21,9,20,5,10,28,26,27,13,15,14,8,4,25,16,11,19,7,2},
	{1,3,8,2,9,5,22,12,20,23,17,7,18,24,21,28,16,19,27,13,26,14,15,10,6,11,4,25},
	{1,21,4,24,17,28,9,16,10,23,3,14,7,18,6,5,13,15,22,11,12,20,26,19,8,2,25,27},
	{2,20,4,8,17,25,16,13,11,18,1,24,14,28,6,23,19,10,26,5,22,21,15,12,3,27,7,9},
	{26,20,2,24,11,21,9,1,19,23,10,4,14,15,12,18,17,7,6,13,5,22,28,25,16,8,27,3}, //145
	{13,10,25,16,12,19,23,18,17,22,2,8,9,3,5,11,14,27,4,26,1,7,28,6,21,24,20,15},
	{4,9,12,3,10,6,5,16,21,26,22,18,28,7,14,11,23,19,1,24,17,2,13,8,15,20,25,27},
	{6,8,18,24,25,21,12,2,1,15,27,23,3,13,16,19,20,28,9,7,22,17,4,5,10,11,14,26},
	{8,9,26,18,23,5,3,16,27,13,14,19,28,11,6,15,2,12,4,17,25,24,22,10,1,20,21,7},
	{4,12,23,17,24,2,18,1,5,14,25,6,11,3,15,28,20,21,26,16,9,22,7,10,13,19,8,27}, //150
	{15,22,20,24,4,23,5,26,12,2,8,14,11,19,9,21,7,17,6,10,3,1,18,28,25,16,13,27},
	{21,10,8,3,1,26,19,4,23,15,7,17,18,13,5,9,28,22,14,25,6,27,2,12,11,24,20,16},
	{24,16,3,7,5,4,18,21,13,17,10,25,19,9,20,8,14,1,15,28,12,6,2,22,26,27,23,11},
	{16,24,14,8,3,19,23,4,26,10,7,21,15,17,6,20,28,13,2,12,1,5,25,27,9,22,18,11},
	{5,8,18,27,7,9,14,20,11,4,1,25,15,13,10,12,24,26,21,23,22,3,6,16,19,28,17,2}, //155
	{27,28,9,14,13,16,15,17,2,24,23,25,1,20,8,22,3,11,10,6,12,5,26,21,18,7,4,19},
	{24,1,11,23,16,5,22,26,7,6,8,13,20,12,15,19,10,14,25,3,17,4,28,2,21,27,18,9},
	{21,13,10,17,15,7,26,23,4,5,22,11,19,12,6,8,27,25,20,2,3,16,9,28,24,18,14,1},
	{22,1,25,2,21,7,17,16,11,9,10,14,26,19,23,13,6,27,20,18,4,28,24,3,12,8,15,5},
	{8,18,22,4,10,16,17,23,9,7,3,12,14,28,13,2,20,5,21,19,27,1,24,11,6,25,15,26}, //160
	{25,6,17,7,22,13,28,24,2,8,26,5,27,12,11,10,14,19,20,9,4,1,21,23,15,18,3,16},
	{27,28,4,23,22,3,5,13,11,12,15,19,2,8,25,14,20,7,18,24,21,10,26,16,9,17,6,1},
	{14,18,27,7,2,17,22,21,28,11,26,1,20,3,4,9,15,8,10,6,13,23,24,19,25,12,16,5},
	{23,24,20,22,13,2,7,12,25,16,10,8,1,6,11,21,27,18,17,28,3,19,9,4,14,15,26,5},
	{5,25,11,28,17,27,22,10,2,20,14,3,15,18,4,24,21,6,23,9,8,19,16,13,7,1,12,26}, //165
	{7,3,18,12,28,15,19,25,24,13,16,9,1,26,23,4,27,6,22,10,2,8,20,17,14,5,11,21},
	{24,20,16,27,5,8,26,15,19,7,2,25,1,22,4,28,13,23,17,14,9,11,10,6,21,18,3,12},
	{24,22,2,20,14,27,18,21,3,5,11,25,6,28,4,16,17,12,13,15,9,26,23,10,8,1,7,19},
	{22,5,3,23,16,13,12,14,19,7,27,18,2,28,17,8,9,26,11,1,6,20,4,21,24,25,15,10},
	{5,16,27,18,11,3,28,13,23,1,2,22,7,26,12,21,20,6,9,17,24,10,4,19,8,25,14,15}, //170
	{17,1,4,15,18,23,19,24,2,13,25,6,12,8,10,21,3,16,7,26,27,22,11,14,20,5,28,9},
	{1,3,11,25,5,22,23,4,10,7,9,27,18,21,2,20,19,12,17,15,16,24,6,28,13,26,8,14},
	{9,6,19,15,7,20,24,14,28,5,17,23,26,16,18,22,12,8,27,13,21,4,3,1,11,2,25,10},
	{2,23,4,3,18,28,11,20,7,9,25,16,22,21,8,14,17,12,5,26,27,24,6,15,13,1,10,19},
	{18,13,28,1,24,11,22,15,14,27,10,16,2,6,8,4,19,5,25,21,9,7,26,17,12,23,3,20}, //175
	{19,3,14,8,7,5,17,6,27,18,15,26,11,21,23,16,9,12,22,2,28,4,20,1,10,13,24,25},
	{11,15,7,20,18,28,25,4,6,21,23,17,26,9,24,1,2,12,10,22,16,13,19,27,8,5,14,3},
	{2,17,20,8,25,28,22,23,21,11,6,10,27,18,9,5,1,16,24,15,12,13,19,26,3,7,14,4},
	{16,9,23,17,7,13,11,8,22,3,6,25,24,2,15,28,1,14,4,12,26,5,10,21,18,27,19,20},
	{5,15,25,24,7,17,1,16,14,13,26,4,9,27,6,3,8,23,19,21,18,2,28,11,20,22,10,12}, //180
	{2,5,25,21,4,13,8,9,28,19,18,23,3,1,12,6,7,14,20,15,10,22,24,26,11,17,16,27},
	{24,25,1,15,21,23,12,13,5,19,7,22,6,11,4,17,27,16,28,18,2,14,8,26,3,9,10,20},
	{25,13,8,5,1,15,10,11,7,27,2,21,16,22,9,17,6,12,24,18,20,4,23,19,28,3,26,14},
	{16,21,3,9,24,1,13,12,25,17,4,18,14,8,11,15,28,20,7,5,19,2,26,10,23,22,6,27},
	{24,8,16,13,19,6,21,14,2,5,4,9,12,27,20,26,15,10,3,22,18,11,23,7,1,17,28,25}, //185
	{2,3,24,20,22,28,4,23,10,27,5,1,9,14,7,12,8,19,16,26,13,17,11,25,18,21,6,15},
	{21,13,15,5,22,12,27,9,8,6,23,20,18,14,24,17,19,25,1,7,10,4,2,26,16,3,11,28},
	{6,4,16,10,27,8,24,25,18,23,13,5,21,2,19,12,15,20,3,9,11,26,1,14,7,22,17,28},
	{3,13,1,26,28,17,18,2,4,14,10,20,12,15,27,21,9,19,16,25,22,6,8,24,7,23,5,11},
	{25,13,10,16,4,3,24,20,1,12,17,18,26,22,14,21,9,11,23,5,7,2,28,6,19,27,15,8}, //190
	{10,27,13,3,5,19,14,1,24,23,4,2,8,25,26,7,11,18,9,21,12,16,28,17,20,22,6,15},
	{8,10,4,1,7,9,11,27,5,18,12,22,19,20,23,26,13,2,14,28,3,17,6,24,21,16,15,25},
	{8,15,14,11,3,20,1,16,18,10,21,9,6,25,2,12,19,22,4,27,28,5,26,23,17,13,24,7},
	{18,25,24,17,7,12,5,6,21,8,3,11,9,23,13,19,14,4,15,26,28,20,27,22,10,16,1,2},
	{28,1,3,21,14,2,20,12,27,19,22,26,15,5,7,13,11,16,23,6,17,24,18,10,4,8,9,25}, //195
	{3,28,7,16,10,19,18,21,23,17,14,11,4,5,20,12,2,22,26,6,13,8,1,25,9,27,24,15},
	{25,7,12,8,18,16,17,4,13,2,26,14,19,3,22,15,11,6,24,28,27,9,21,20,23,5,10,1},
	{15,6,14,4,20,17,18,22,23,19,12,5,16,27,25,3,2,8,11,21,28,1,26,24,10,9,13,7},
	{2,10,21,12,17,22,3,1,23,24,25,11,15,19,7,28,13,27,20,8,16,14,6,5,4,9,26,18},
	{4,14,25,2,18,11,13,12,10,6,26,7,5,23,27,1,19,21,24,15,3,17,16,20,28,8,22,9}, //200
	{10,20,5,4,1,25,7,27,21,26,17,9,22,2,23,16,15,13,12,14,3,24,11,28,18,8,6,19},
	{6,22,7,23,24,13,9,11,25,15,27,19,2,3,20,16,10,1,8,26,17,14,12,18,21,28,5,4},
	{23,22,9,3,11,8,15,5,25,10,28,1,21,19,13,27,4,2,17,12,14,7,18,26,6,24,20,16},
	{4,3,13,11,9,20,24,17,21,14,19,6,1,10,15,28,2,26,12,22,25,27,7,18,23,16,8,5},
	{18,19,24,21,8,20,22,13,27,2,11,4,25,3,26,5,14,23,12,15,6,7,1,10,28,16,9,17}, //205
	{2,20,23,10,12,13,15,9,16,18,4,17,6,19,21,7,27,28,8,11,14,26,5,22,25,1,3,24},
	{3,15,16,14,8,24,10,13,28,21,2,23,18,11,6,25,5,20,4,17,26,22,9,12,1,27,7,19},
	{2,17,1,19,27,25,16,21,11,20,4,13,7,9,14,28,24,5,23,12,18,8,3,26,10,15,6,22},
	{16,1,13,15,20,23,9,21,5,14,12,10,2,25,6,28,27,18,17,3,19,11,26,24,7,8,22,4},
	{14,20,24,1,27,25,19,4,16,7,21,2,23,17,8,26,6,13,3,18,5,9,11,15,10,12,22,28}, //210
	{13,21,6,24,4,7,19,2,23,28,15,16,8,17,27,22,25,14,20,12,5,9,3,10,18,1,26,11},
	{26,11,7,13,9,4,22,21,28,3,24,23,1,17,12,15,25,18,20,10,16,27,14,6,2,8,19,5},
	{19,18,11,13,24,8,28,16,10,9,23,14,6,27,7,26,25,12,21,4,17,5,3,20,22,2,15,1},
	{7,24,23,19,28,22,12,21,20,4,10,18,17,15,25,8,11,6,3,1,5,26,13,9,14,16,27,2},
	{5,14,12,24,27,2,23,20,18,16,9,1,6,19,22,10,11,13,8,7,17,25,21,26,4,15,3,28}, //215
	{20,22,1,17,25,10,27,26,13,12,2,8,6,23,9,4,14,24,7,5,3,21,19,18,16,15,28,11},
	{3,7,17,4,28,19,8,24,1,20,21,6,2,13,5,18,11,10,25,22,15,14,27,26,9,12,23,16},
	{26,10,25,17,23,3,22,6,2,21,7,4,1,14,12,5,9,15,28,18,27,20,8,11,16,24,19,13},
	{11,25,2,24,13,21,12,22,4,26,19,5,20,15,6,16,8,10,23,7,17,28,18,9,14,3,1,27},
	{4,19,3,18,20,9,15,8,28,6,24,27,11,14,7,12,10,17,21,5,13,25,1,23,16,22,2,26}, //220
	{24,15,17,5,26,10,21,3,16,18,23,2,27,12,22,19,13,20,6,25,1,4,8,11,28,9,14,7},
	{8,15,1,21,12,7,11,25,14,27,9,16,5,13,17,19,10,24,28,4,26,2,23,6,3,18,22,20},
	{23,16,5,22,11,19,28,14,4,6,10,24,27,18,15,3,26,8,7,20,9,1,2,21,25,13,12,17},
	{27,4,25,19,26,6,15,8,16,18,2,20,1,10,28,3,5,11,23,17,12,21,9,14,24,7,13,22},
	{28,9,12,21,17,18,22,24,6,20,11,5,7,15,25,27,8,3,13,10,23,1,19,2,16,4,14,26}, //225
	{1,25,20,13,21,16,28,4,8,24,3,11,10,26,6,27,23,15,2,19,12,17,18,5,9,22,7,14},
	{12,9,7,5,25,6,17,8,15,14,22,26,10,20,16,24,4,19,13,27,3,1,28,21,18,23,11,2},
	{17,23,14,28,10,21,4,19,18,25,5,26,27,16,8,1,7,11,22,3,9,20,12,2,24,13,15,6},
	{16,1,9,10,13,12,27,18,26,2,22,14,7,11,6,5,4,15,28,20,25,21,17,3,8,19,24,23},
	{4,1,22,26,6,21,16,25,3,9,24,19,27,18,20,17,15,28,11,14,8,12,13,2,7,5,23,10}, //230
	{28,27,17,11,9,5,1,22,18,24,19,25,14,12,16,21,26,3,20,4,23,2,15,6,7,13,10,8},
	{21,12,6,22,4,18,8,24,16,17,14,20,25,13,9,15,27,28,5,11,23,10,7,3,26,1,19,2},
	{14,28,19,3,20,2,16,25,27,8,12,26,22,1,18,7,13,9,6,15,4,10,23,21,11,17,24,5},
	{23,1,19,17,5,6,14,20,24,12,11,18,21,13,25,4,8,28,15,10,3,27,2,22,7,9,16,26},
	{1,3,16,15,20,2,25,7,27,14,28,26,12,23,6,4,18,19,5,9,24,11,13,22,8,17,21,10}, //235
	{23,7,2,19,17,1,21,20,26,15,5,8,3,12,22,28,27,24,16,25,13,18,6,9,4,10,14,11},
	{12,26,5,20,7,2,11,16,19,8,22,21,9,10,13,23,4,25,1,27,17,15,6,24,18,14,28,3},
	{1,27,21,12,22,5,18,14,3,17,10,2,26,20,24,25,28,7,23,13,8,6,4,16,19,11,9,15},
	{21,4,3,24,19,6,23,26,7,13,2,22,5,10,14,20,27,11,16,15,8,18,12,28,25,17,9,1},
	{14,22,8,5,21,23,3,13,19,2,15,9,20,12,16,27,25,17,11,26,4,7,24,6,10,28,1,18}, //240
	{6,23,13,19,11,24,27,16,9,14,8,25,4,12,15,5,17,3,28,2,21,20,22,10,1,18,26,7},
	{28,9,17,7,10,23,20,16,21,25,8,1,14,26,13,24,5,19,4,11,22,15,2,12,27,6,3,18},
	{1,21,10,24,6,16,26,5,25,7,17,20,14,9,3,18,13,28,12,4,8,11,19,15,23,22,2,27},
	{1,8,17,24,23,21,18,19,4,26,7,28,14,2,10,25,16,15,13,22,27,5,6,20,11,12,9,3},
	{20,10,23,24,17,6,21,27,5,9,7,8,12,25,2,11,15,13,18,4,26,16,22,3,19,14,1,28}, //245
	{4,17,19,25,14,27,5,24,12,26,18,2,20,21,23,1,10,15,16,3,13,28,8,11,22,9,7,6},
	{10,19,12,24,18,4,8,15,6,17,1,22,9,21,13,25,26,20,3,14,23,7,2,16,11,28,5,27},
	{6,9,3,20,26,25,18,7,28,11,10,19,17,2,12,21,4,15,8,22,1,5,16,23,14,13,24,27},
	{25,9,24,28,7,18,22,21,1,19,12,20,3,15,10,4,16,8,5,13,17,14,6,26,11,23,2,27},
	{28,5,26,3,19,13,21,10,4,22,23,24,12,1,16,18,8,15,14,7,25,17,27,9,2,11,20,6}, //250
	{5,20,18,8,6,26,13,11,28,24,9,2,4,17,15,22,21,19,1,12,16,3,23,27,14,10,25,7},
	{14,20,27,28,25,21,6,18,1,13,10,16,4,15,24,19,11,3,22,26,2,7,17,8,5,23,12,9},
	{21,26,16,20,2,12,19,22,5,7,1,4,6,24,23,15,14,9,10,25,28,18,27,3,8,11,17,13},
	{17,25,15,3,12,14,9,10,16,21,7,5,8,24,22,13,6,4,11,23,26,19,28,20,27,2,1,18},
	{3,8,19,14,7,21,10,28,9,11,26,12,20,2,23,15,4,16,13,22,1,18,24,6,27,17,25,5}} //255
