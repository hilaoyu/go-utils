package utilOrm

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
	Timezone   = "Asia/Shanghai"
)

type OrmTime time.Time

func ormTimeLoc() *time.Location {
	loc, err1 := time.LoadLocation(Timezone)
	if nil != err1 {
		loc = time.Local
	}

	return loc
}

func (t OrmTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).In(ormTimeLoc()).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *OrmTime) UnmarshalJSON(data []byte) (err error) {

	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), ormTimeLoc())
	*t = OrmTime(now)
	return
}

func (t OrmTime) String() string {
	return time.Time(t).In(ormTimeLoc()).Format(TimeFormat)
}

func (t OrmTime) local() time.Time {
	return time.Time(t).In(ormTimeLoc())
}

func (t OrmTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t).In(ormTimeLoc())
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *OrmTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = OrmTime(value.In(ormTimeLoc()))
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
