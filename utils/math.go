package utils

func CeilInt(a, b int) int {
	if b <= 0 {
		return a
	}
	return (a + b - 1) / b
}

func CeilUint8(a, b uint8) uint8 {
	if b <= 0 {
		return a
	}
	return (a + b - 1) / b
}
func CeilUint32(a, b uint32) uint32 {
	if b <= 0 {
		return a
	}
	return (a + b - 1) / b
}

func CeilUint64(a, b uint64) uint64 {
	if b <= 0 {
		return a
	}
	return (a + b - 1) / b
}
