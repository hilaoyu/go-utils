package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

func ReTry(callback func() bool, times int, step time.Duration) {
	for i := 0; i < times; i++ {
		if callback() {
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
