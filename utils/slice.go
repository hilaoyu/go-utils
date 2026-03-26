package utils

import (
	"reflect"
	"slices"
	"strings"

	"github.com/hilaoyu/go-utils/utilRandom"
)

// SliceDifference 切片差集
func SliceDifference[T comparable](a, b []T) []T {
	setB := make(map[T]struct{}, len(b))
	for _, v := range b {
		setB[v] = struct{}{}
	}

	seen := make(map[T]struct{})
	res := make([]T, 0, len(a))

	for _, v := range a {
		if _, inB := setB[v]; !inB {
			if _, exists := seen[v]; !exists {
				res = append(res, v)
				seen[v] = struct{}{}
			}
		}
	}
	return res
}

// SliceShift 返回并删除切片第一个元素
func SliceShift[E any](s *[]E) (E) {
	if len(*s) == 0 {
		var zero E
		return zero
	}

	e := (*s)[0]

	// 避免内存泄漏（关键）
	var zero E
	(*s)[0] = zero

	*s = (*s)[1:]

	return e
}


// SliceUnShift 在切片头部加入一个元素
func SliceUnShift[E any](s []E, e E) []E {
	s = append(s, e)
	copy(s[1:], s[:len(s)-1])
	s[0] = e
	return s
}

// SlicePop 返回并删除切片最后一个元素
func SlicePop[E any](s *[]E) (E) {
	if len(*s) == 0 {
		var zero E
		return zero
	}

	last := len(*s) - 1
	e := (*s)[last]
	*s = (*s)[:last]

	return e
}

// SlicePush 在切片尾部加入一个元素
func SlicePush[ E any](s []E, e E) {
	s = append(s, e)
	return
}

// SliceRandom 切片随机返回一个元素
func SliceRandom[E any](s []E) (e E) {
	randomIndex := utilRandom.RandInt(len(s))
	e = s[randomIndex]
	return
}

// SliceFind 切片查找元素
func SliceFind[ E any](s []E, f func(E) bool) (e E) {
	i := slices.IndexFunc(s, f)
	if i >= 0 {
		e = s[i]
	}
	return
}

// SliceContains 切片是否包含某个元素
func SliceContains[ E any](s []E, e E) bool {
	i := slices.IndexFunc(s, func(f E) bool {
		return reflect.DeepEqual(f, e)
	})
	return i >= 0
}

// SliceFilter 切片筛选
func SliceFilter[E any](s []E, f func(E) bool) (s1 []E) {
	for _, e := range s {
		if f(e) {
			s1 = append(s1, e)
		}
	}
	return
}

// SliceDelElement 切片删除元素
func SliceDelElement[E any](s *[]E, f func(E) bool) {
	s1 := *s
	*s = nil
	for _, e := range s1 {
		if !f(e) { *s = append(*s, e) }
	}
	return
	
}

// SliceShuffle 切片乱序
func SliceShuffle[E any](s []E) []E {
	mi := len(s) - 1
	for i := range s {
		j := utilRandom.RandInt(mi) // 生成一个 [0, mi) 范围内的随机索引
		s[i], s[j] = s[j], s[i]     // 交换元素
	}
	return s
}

// SliceStringFilterUniqueNonEmpty 字符串切片去空和去重
func SliceStringFilterUniqueNonEmpty(ss []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}
