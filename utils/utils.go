package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"time"
)

// EachMapSort /*
// 以map的key(int\float\string)排序遍历map
// eachMap      ->  待遍历的map
// eachFunc     ->  map遍历接收，入参应该符合map的key和value
// 需要对传入类型进行检查，不符合则直接panic提醒进行代码调整

func EachMapSort(eachMap interface{}, eachFunc interface{}) error {
	eachMapValue := reflect.ValueOf(eachMap)
	eachFuncValue := reflect.ValueOf(eachFunc)
	eachMapType := eachMapValue.Type()
	eachFuncType := eachFuncValue.Type()
	if eachMapValue.Kind() != reflect.Map {
		return errors.New("ksort.EachMap failed. parameter \"eachMap\" type must is map[...]...{}")
	}
	if eachFuncValue.Kind() != reflect.Func {
		return errors.New("ksort.EachMap failed. parameter \"eachFunc\" type must is func(key ..., value ...)")
	}
	if eachFuncType.NumIn() != 2 {
		return errors.New("ksort.EachMap failed. \"eachFunc\" input parameter count must is 2")
	}
	if eachFuncType.In(0).Kind() != eachMapType.Key().Kind() {
		return errors.New("ksort.EachMap failed. \"eachFunc\" input parameter 1 type not equal of \"eachMap\" key")
	}
	if eachFuncType.In(1).Kind() != eachMapType.Elem().Kind() {
		return errors.New("ksort.EachMap failed. \"eachFunc\" input parameter 2 type not equal of \"eachMap\" value")
	}

	// 对key进行排序
	// 获取排序后map的key和value，作为参数调用eachFunc即可
	switch eachMapType.Key().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		keys := make([]int, 0)
		keysMap := map[int]reflect.Value{}
		for _, value := range eachMapValue.MapKeys() {
			keys = append(keys, int(value.Int()))
			keysMap[int(value.Int())] = value
		}
		sort.Ints(keys)
		for _, key := range keys {
			eachFuncValue.Call([]reflect.Value{keysMap[key], eachMapValue.MapIndex(keysMap[key])})
		}
	case reflect.Float64, reflect.Float32:
		keys := make([]float64, 0)
		keysMap := map[float64]reflect.Value{}
		for _, value := range eachMapValue.MapKeys() {
			keys = append(keys, float64(value.Float()))
			keysMap[float64(value.Float())] = value
		}
		sort.Float64s(keys)
		for _, key := range keys {
			eachFuncValue.Call([]reflect.Value{keysMap[key], eachMapValue.MapIndex(keysMap[key])})
		}
	case reflect.String:
		keys := make([]string, 0)
		keysMap := map[string]reflect.Value{}
		for _, value := range eachMapValue.MapKeys() {
			keys = append(keys, value.String())
			keysMap[value.String()] = value
		}
		sort.Strings(keys)
		for _, key := range keys {
			eachFuncValue.Call([]reflect.Value{keysMap[key], eachMapValue.MapIndex(keysMap[key])})
		}
	default:
		return errors.New("\"eachMap\" key type must is int or float or string")
	}
	return nil
}

// ReTry /** time: 0 forever

func ReTry(callback func() bool, times int, step time.Duration) {
	forever := false
	if times <= 0 {
		forever = true
		times = 1
	}
	for {
		if callback() {
			return
		}
		if !forever {
			times--
		}
		if times <= 0 {
			return
		}
		time.Sleep(step)
	}

}

func InterfaceToStruct(out interface{}, in interface{}) error {
	jsonStr, err := json.Marshal(in)
	if nil != err {
		return err
	}

	err = json.Unmarshal(jsonStr, &out)

	return err

}

func GetInterfaceFiledValue(v interface{}, fieldName string) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(v)).FieldByName(fieldName)
}

func MakeInstanceFromSlice(v interface{}) (i interface{}, err error) {
	typ := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Slice {
		return nil, fmt.Errorf("argument must be a slice")
	}

	i = reflect.MakeSlice(typ.Type(), 1, 1).Index(0).Interface()

	return
}

func MapStringValues(m map[string]string) (v []string) {
	for _, s := range m {
		v = append(v, s)
	}
	return
}

func ReverseBits(b byte) byte {
	var reverse = [256]int{
		0, 128, 64, 192, 32, 160, 96, 224,
		16, 144, 80, 208, 48, 176, 112, 240,
		8, 136, 72, 200, 40, 168, 104, 232,
		24, 152, 88, 216, 56, 184, 120, 248,
		4, 132, 68, 196, 36, 164, 100, 228,
		20, 148, 84, 212, 52, 180, 116, 244,
		12, 140, 76, 204, 44, 172, 108, 236,
		28, 156, 92, 220, 60, 188, 124, 252,
		2, 130, 66, 194, 34, 162, 98, 226,
		18, 146, 82, 210, 50, 178, 114, 242,
		10, 138, 74, 202, 42, 170, 106, 234,
		26, 154, 90, 218, 58, 186, 122, 250,
		6, 134, 70, 198, 38, 166, 102, 230,
		22, 150, 86, 214, 54, 182, 118, 246,
		14, 142, 78, 206, 46, 174, 110, 238,
		30, 158, 94, 222, 62, 190, 126, 254,
		1, 129, 65, 193, 33, 161, 97, 225,
		17, 145, 81, 209, 49, 177, 113, 241,
		9, 137, 73, 201, 41, 169, 105, 233,
		25, 153, 89, 217, 57, 185, 121, 249,
		5, 133, 69, 197, 37, 165, 101, 229,
		21, 149, 85, 213, 53, 181, 117, 245,
		13, 141, 77, 205, 45, 173, 109, 237,
		29, 157, 93, 221, 61, 189, 125, 253,
		3, 131, 67, 195, 35, 163, 99, 227,
		19, 147, 83, 211, 51, 179, 115, 243,
		11, 139, 75, 203, 43, 171, 107, 235,
		27, 155, 91, 219, 59, 187, 123, 251,
		7, 135, 71, 199, 39, 167, 103, 231,
		23, 151, 87, 215, 55, 183, 119, 247,
		15, 143, 79, 207, 47, 175, 111, 239,
		31, 159, 95, 223, 63, 191, 127, 255,
	}

	return byte(reverse[int(b)])
}

func SliceShift[S ~*[]E, E any](s S) (e E) {
	if len(*s) <= 0 {
		return
	}
	e = (*s)[0]
	*s = (*s)[1:]
	return
}
func SlicePop[S ~*[]E, E any](s S) (e E) {
	sl := len(*s)
	if sl <= 0 {
		return
	}
	e = (*s)[sl-1]
	*s = (*s)[:sl-1]
	return
}

func SliceFind[S ~[]E, E any](s S, f func(E) bool) (e E) {
	i := slices.IndexFunc(s, f)
	if i >= 0 {
		e = s[i]
	}
	return
}
func SliceContains[S ~[]E, E any](s S, e E) bool {
	i := slices.IndexFunc(s, func(f E) bool {
		return reflect.DeepEqual(f, e)
	})
	return i >= 0
}
