package checkin

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type CheckIn interface {
	// 用户签到
	UserCheckIn()
	// 某个用户本月的签到记录
	UserCheckInRecord()
	// 某个用户本月签到总数
	UserCheckInCount()
	// 某个用户是否已签到
	IsCheckIn()
	// 某天的总签到用户数
	DayRecord()
}

type CheckInStruct struct {
	Conn *redis.Client
}

var (
	CheckInDayKey   = "user:checkin:day:"
	CheckInMonthKey = "user:checkin:month:"
	CheckInCountKey = "user:checkin:count:"
)

var (
	ctx                = context.Background()
	year, month, today = time.Now().Date()
)

var (
	todayCheckInDayKey = CheckInDayKey + strconv.Itoa(today)
	monthCheckInKey    = CheckInMonthKey + strconv.Itoa(int(month))
)

// 用户签到
func (this *CheckInStruct) UserCheckIn(userId int) {

	this.Conn.SetBit(ctx, todayCheckInDayKey, 1, userId)
	this.Conn.SetBit(ctx, monthCheckInKey, 1, userId)
}

// todo 某个用户本月的签到记录
func (this *CheckInStruct) UserCheckInRecord(userId int) string {
	field := this.Conn.BitField(ctx, monthCheckInKey, "get", "u28", 0)
	return field.String()
}

// 某个用户本月签到总数
func (this *CheckInStruct) UserCheckInCount() {

}

// 某天某个用户是否已签到
func (this *CheckInStruct) IsCheckIn(userId int, day string) string {
	dayCheckInDayKey := CheckInDayKey + day
	bit := this.Conn.GetBit(ctx, dayCheckInDayKey, int64(userId))

	return bit.String()
}

func (this *CheckInStruct) DayRecord() {

}
