package utils

import (
	"github.com/hilaoyu/go-utils/utilRandom"
	"reflect"
	"slices"
)

// SliceShift 返回并删除切片第一个元素
func SliceShift[S ~*[]E, E any](s S) (e E) {
	if len(*s) <= 0 {
		return
	}
	e = (*s)[0]
	*s = (*s)[1:]
	return
}

// SliceUnShift 在切片头部加入一个元素
func SliceUnShift[S ~*[]E, E any](s S, e E) {
	*s = append([]E{e}, *s...)
	return
}

// SlicePop 返回并删除切片最后一个元素
func SlicePop[S ~*[]E, E any](s S) (e E) {
	sl := len(*s)
	if sl <= 0 {
		return
	}
	e = (*s)[sl-1]
	*s = (*s)[:sl-1]
	return
}

// SlicePush 在切片尾部加入一个元素
func SlicePush[S ~*[]E, E any](s S, e E) {
	*s = append(*s, e)
	return
}

// SliceRandom 切片随机返回一个元素
func SliceRandom[S ~[]E, E any](s S) (e E) {
	randomIndex := utilRandom.RandInt(len(s))
	e = s[randomIndex]
	return
}

// SliceFind 切片查找元素
func SliceFind[S ~[]E, E any](s S, f func(E) bool) (e E) {
	i := slices.IndexFunc(s, f)
	if i >= 0 {
		e = s[i]
	}
	return
}

// SliceContains 切片是否包含某个元素
func SliceContains[S ~[]E, E any](s S, e E) bool {
	i := slices.IndexFunc(s, func(f E) bool {
		return reflect.DeepEqual(f, e)
	})
	return i >= 0
}

// SliceFilter 切片筛选
func SliceFilter[S ~[]E, E any](s S, f func(E) bool) (s1 S) {
	for _, e := range s {
		if f(e) {
			s1 = append(s1, e)
		}
	}
	return
}

// SliceDelElement 切片删除元素
func SliceDelElement[S ~*[]E, E any](s S, f func(E) bool) {
	s1 := *s
	*s = nil
	for _, e := range s1 {
		if !f(e) {
			*s = append(*s, e)
		}
	}
	return
}

// SliceShuffle 切片乱序
func SliceShuffle[S ~[]E, E any](s S) S {
	mi := len(s) - 1
	for i := range s {
		j := utilRandom.RandInt(mi) // 生成一个 [0, mi) 范围内的随机索引
		s[i], s[j] = s[j], s[i]     // 交换元素
	}
	return s
}
