package utilStr

import (
	"strconv"
	"strings"
)

func SubString(source string, start int, length int) string {
	var r = []rune(source)
	realLength := len(r)
	var end int = realLength
	if start < 0 {
		start = realLength + start
	}

	if length > 0 {
		end = start + length
	}

	if end > realLength {
		end = realLength
	}
	if start < 0 || start >= realLength {
		return ""
	}

	if start == 0 && end == realLength {
		return source
	}

	if start == end {
		return string(r[start])
	}

	var substring = ""
	for i := start; i < realLength; i++ {

		if i >= end {
			break
		}
		substring += string(r[i])
	}

	return substring
}

func After(haystack string, needle string) string {
	temp := haystack
	arr := strings.SplitN(haystack, needle, 2)
	if len(arr) >= 2 {
		temp = arr[1]
	}
	return temp
}
func AfterLast(haystack string, needle string) string {

	arr := strings.Split(haystack, needle)

	return arr[len(arr)-1]
}

func Before(haystack string, needle string) string {

	arr := strings.Split(haystack, needle)
	arrLen := len(arr)
	if arrLen >= 2 {
		arr = arr[0 : arrLen-1]
	}
	temp := ""
	for i, v := range arr {
		if i > 0 {
			v = needle + v
		}
		temp += v
	}

	return temp
}
func BeforeFirst(haystack string, needle string) string {

	arr := strings.SplitN(haystack, needle, 2)

	return arr[0]
}

func IsDigit(s string) bool {
	_, err := strconv.ParseFloat(s, 64)

	return nil == err
}
