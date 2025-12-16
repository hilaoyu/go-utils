package utilOrm

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
	Timezone   = "Asia/Shanghai"
)

var parseFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	time.RFC3339,
	time.RFC3339Nano,
}

type OrmTime struct {
	time.Time
}

func NewOrmTime(timeTime time.Time) *OrmTime {
	return &OrmTime{timeTime}
}

func OrmTimeParse(value string, layout ...string) (t *OrmTime, err error) {
	s := strings.TrimSpace(value)
	if s == "" {
		return nil, fmt.Errorf("value is empty")
	}
	ts, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		t = NewOrmTime(time.Unix(ts, 0))
		return
	}
	err = nil

	if len(layout) <= 0 {
		layout = parseFormats
	}

	// try all formats
	for _, f := range layout {
		if to, err1 := time.ParseInLocation(f, s, ormTimeLoc()); err1 == nil {
			t = NewOrmTime(to)
			break
		}
	}

	if t == nil {
		err = fmt.Errorf("invalid time format: %s", value)
	}

	return
}

func ormTimeLoc() *time.Location {
	loc, err1 := time.LoadLocation(Timezone)
	if nil != err1 {
		loc = time.Local
	}

	return loc
}

func (t OrmTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("\"\""), nil
	}
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = t.In(ormTimeLoc()).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *OrmTime) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(string(data), "\"")
	if nil == data || len(data) == 0 || "null" == strings.ToLower(s) {
		return nil
	}

	return t.UnmarshalText([]byte(s))

}

func (t *OrmTime) UnmarshalText(text []byte) error {
	tmp, err := OrmTimeParse(string(text), Timezone)
	if nil != err {
		return err
	}
	t.Time = tmp.Time

	return nil
}

func (t *OrmTime) String() string {
	if nil == t || t.IsZero() {
		return ""
	}
	return t.In(ormTimeLoc()).Format(TimeFormat)
}

func (t *OrmTime) local() time.Time {
	return t.In(ormTimeLoc())
}

func (t OrmTime) Value() (driver.Value, error) {

	if t.IsZero() {
		return nil, nil
	}
	var ti = t.In(ormTimeLoc())

	return ti, nil
}

func (t *OrmTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		//*t = OrmTime(value.In(ormTimeLoc()))
		t.Time = value.In(ormTimeLoc())
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
