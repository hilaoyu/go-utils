package utilTime

import (
	"time"
)

var timeFormat = "2006-01-02 15:04:05"
var timeLocation = "Asia/Shanghai"

func SetTimeFormat(v string) {
	timeFormat = v
}
func GetTimeFormat() string {
	return timeFormat
}
func SetTimeLocation(v string) {
	timeLocation = v
}
func GetTimeLocation() string {
	return timeLocation
}

func TimeGetFormatLayout(format ...string) string {
	tempTimeFormat := timeFormat
	if len(format) > 0 && "" != format[0] {
		tempTimeFormat = format[0]
	}
	return tempTimeFormat
}
func TimeGetLocation(location ...string) (*time.Location, error) {
	tempLocation := timeLocation
	if len(location) > 0 && "" != location[0] {
		tempLocation = location[0]
	}
	return time.LoadLocation(tempLocation)
}

func TimeFormat(t time.Time, format ...string) string {
	loc, _ := TimeGetLocation()
	return t.In(loc).Format(TimeGetFormatLayout(format...))
}

func TimeParse(s string, format ...string) time.Time {
	loc, _ := TimeGetLocation()

	t, err := time.ParseInLocation(TimeGetFormatLayout(format...), s, loc)
	if nil != err {
		//fmt.Println(TimeGetFormatLayout(format...), err)
		t = time.Time{}
	}
	return t
}

func TimeNow() time.Time {
	loc, _ := TimeGetLocation()
	return time.Now().In(loc)
}
func TimeNowString(format ...string) string {
	return TimeFormat(TimeNow(), format...)
}
