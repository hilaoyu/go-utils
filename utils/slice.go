package utils

import (
	"github.com/hilaoyu/go-utils/utilRandom"
	"reflect"
	"slices"
)

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
func SliceRandom[S ~[]E, E any](s S) (e E) {
	randomIndex := utilRandom.RandInt(len(s))
	e = s[randomIndex]
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

func SliceFilter[S ~[]E, E any](s S, f func(E) bool) (s1 S) {
	for _, e := range s {
		if f(e) {
			s1 = append(s1, e)
		}
	}
	return
}

func SliceShuffle[S ~[]E, E any](s S) S {
	mi := len(s) - 1
	for i := range s {
		j := utilRandom.RandInt(mi) // 生成一个 [0, mi) 范围内的随机索引
		s[i], s[j] = s[j], s[i]     // 交换元素
	}
	return s
}
