package utilExcel

import "strings"

var (
	columnIndexList = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func ColumnIndexToTitle(num int) string {
	num = num - 1
	var column = columnIndexList[num%26]
	for num = num / 26; num > 0; num = num/26 - 1 {
		column = columnIndexList[(num-1)%26] + column
	}
	return column
}

func ColumnTitleToIndex(columnTitle string) int {
	columnTitle = strings.ToUpper(columnTitle)
	const cntLetter = 26
	var result int

	n := 1
	for i := len(columnTitle) - 1; i >= 0; i-- {
		// 字母转换成对应的 26 进制的数字
		// A : 1 , B:2 , Z:26
		val := (int(columnTitle[i]) - 65 + 1) * n
		// 各位结果相加
		result = result + val

		// 从低位到高位，是 26 的 n 次方
		// 最右端：26 的 0 次方
		// 右边第二位：26 的 1 次方
		// 右边第三位：26 的 2 次方
		n = n * cntLetter
	}

	return result
}
