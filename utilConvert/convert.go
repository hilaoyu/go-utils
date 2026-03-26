package utilConvert

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// convert string to specify type

type StrTo string

func (f *StrTo) Set(v string) {
	if v != "" {
		*f = StrTo(v)
	} else {
		f.Clear()
	}
}

func (f *StrTo) Clear() {
	*f = StrTo(0x1E)
}

func (f StrTo) Exist() bool {
	return string(f) != string(0x1E)
}

func (f StrTo) Bool() (bool, error) {
	if f == "on" {
		return true, nil
	}
	return strconv.ParseBool(f.String())
}

func (f StrTo) Float32() (float32, error) {
	v, err := strconv.ParseFloat(f.String(), 32)
	return float32(v), err
}

func (f StrTo) Float64() (float64, error) {
	return strconv.ParseFloat(f.String(), 64)
}

func (f StrTo) Int() (int, error) {
	v, err := strconv.ParseInt(f.String(), 10, 32)
	return int(v), err
}

func (f StrTo) Int8() (int8, error) {
	v, err := strconv.ParseInt(f.String(), 10, 8)
	return int8(v), err
}

func (f StrTo) Int16() (int16, error) {
	v, err := strconv.ParseInt(f.String(), 10, 16)
	return int16(v), err
}

func (f StrTo) Int32() (int32, error) {
	v, err := strconv.ParseInt(f.String(), 10, 32)
	return int32(v), err
}

func (f StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(f.String(), 10, 64)
	return int64(v), err
}

func (f StrTo) Uint() (uint, error) {
	v, err := strconv.ParseUint(f.String(), 10, 32)
	return uint(v), err
}

func (f StrTo) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v), err
}

func (f StrTo) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(f.String(), 10, 16)
	return uint16(v), err
}

func (f StrTo) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(f.String(), 10, 32)
	return uint32(v), err
}

func (f StrTo) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(f.String(), 10, 64)
	return uint64(v), err
}

func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}

// convert any type to string
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}


// convert any numeric value to int64
func ToInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	case string:
		d, err = strconv.ParseInt(val.String(), 10, 64)

	case float64, float32:
		d, err = strconv.ParseInt(fmt.Sprintf("%1.0f", val.Float()), 10, 64)
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}
func ToInt64IgnoreError(value interface{}) (int64){
	v,_ := ToInt64(value)
	return v
}
func ToInt(value interface{}) (d int, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = int(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = int(val.Uint())
	case string:
		d, err = strconv.Atoi(val.String())

	case float64, float32:
		d, err = strconv.Atoi(fmt.Sprintf("%1.0f", val.Float()))

	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}
func ToIntIgnoreError(value interface{}) (int){
	v,_ := ToInt(value)
	return v
}

func ToUint8(value interface{}) (uint8, error) {
	switch v := value.(type) {

	case uint8:
		return v, nil

	case uint, uint16, uint32, uint64:
		n := reflect.ValueOf(v).Uint()
		if n > math.MaxUint8 {
			return 0, fmt.Errorf("value %d overflows uint8", n)
		}
		return uint8(n), nil

	case int, int8, int16, int32, int64:
		n := reflect.ValueOf(v).Int()
		if n < 0 || n > math.MaxUint8 {
			return 0, fmt.Errorf("value %d overflows uint8", n)
		}
		return uint8(n), nil

	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f < 0 || f > math.MaxUint8 {
			return 0, fmt.Errorf("value %f overflows uint8", f)
		}
		return uint8(f), nil

	case string:
		s := strings.TrimSpace(v)
		n, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return 0, err
		}
		return uint8(n), nil

	default:
		return 0, fmt.Errorf("ToUint8: unsupported type %T", value)
	}
}
func ToUint8IgnoreError(value interface{}) (uint8){
	v,_ := ToUint8(value)
	return v
}
func ToUint32(value interface{}) (uint32, error) {
	switch v := value.(type) {

	case uint32:
		return v, nil

	case uint, uint8, uint16, uint64:
		n := reflect.ValueOf(v).Uint()
		if n > math.MaxUint32 {
			return 0, fmt.Errorf("value %d overflows uint32", n)
		}
		return uint32(n), nil

	case int, int8, int16, int32, int64:
		n := reflect.ValueOf(v).Int()
		if n < 0 || n > math.MaxUint32 {
			return 0, fmt.Errorf("value %d overflows uint32", n)
		}
		return uint32(n), nil

	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f < 0 || f > math.MaxUint32 {
			return 0, fmt.Errorf("value %f overflows uint32", f)
		}
		return uint32(f), nil

	case string:
		s := strings.TrimSpace(v)
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(n), nil

	default:
		return 0, fmt.Errorf("ToUint32: unsupported type %T", value)
	}
}
func ToUint32IgnoreError(value interface{}) (uint32){
	v,_ := ToUint32(value)
	return v
}
func ToUint64(value interface{}) (uint64, error) {
	switch v := value.(type) {

	case uint64:
		return v, nil

	case uint, uint8, uint16, uint32:
		return reflect.ValueOf(v).Uint(), nil

	case int, int8, int16, int32, int64:
		n := reflect.ValueOf(v).Int()
		if n < 0 {
			return 0, fmt.Errorf("value %d overflows uint64", n)
		}
		return uint64(n), nil

	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f < 0 || f > math.MaxUint64 {
			return 0, fmt.Errorf("value %f overflows uint64", f)
		}
		return uint64(f), nil

	case string:
		s := strings.TrimSpace(v)
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil

	default:
		return 0, fmt.Errorf("ToUint64: unsupported type %T", value)
	}
}
func ToUint64IgnoreError(value interface{}) (uint64){
	v,_ := ToUint64(value)
	return v
}

type argString []string

func (a argString) Get(i int, args ...string) (r string) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

type argInt []int

func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

type argAny []interface{}

func (a argAny) Get(i int, args ...interface{}) (r interface{}) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}
