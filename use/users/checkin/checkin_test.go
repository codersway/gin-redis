package checkin

import (
	"redis-in-action/conf"
	"strconv"
	"testing"
)

var (
	conn   = conf.Conn()
	client = CheckInStruct{
		Conn: conn,
	}
)

func TestCheckInStruct_UserCheckIn(t *testing.T) {
	client.UserCheckIn(3)
}

func TestCheckInStruct_UserCheckInRecord(t *testing.T) {
	record := client.UserCheckInRecord(1)
	t.Log(record)
}

func TestCheckInStruct_DayRecord(t *testing.T) {

}

func TestCheckInStruct_IsCheckIn(t *testing.T) {
	in := client.IsCheckIn(1, strconv.Itoa(today))
	t.Log(in)
}

func TestCheckInStruct_UserCheckInCount(t *testing.T) {

}
