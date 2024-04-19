package utilRandom

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandInt(n int64) int64 {
	i := rand.Int63n(n)
	return i
}

func RandString(n int, chars ...string) string {
	//letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterBytes := "abcdefhkmnprstuvwxyzABCEFGHIJKLMNPRSTUVWXYZ1234567890"
	if len(chars) >= 1 {
		letterBytes = strings.Join(chars, "")
	}
	b := make([]byte, n)
	letterBytesLen := len(letterBytes)
	for i := range b {
		if 0 == i && len(chars) <= 0 {
			b[i] = letterBytes[rand.Intn(letterBytesLen-10)] //首位为字母
		} else {
			b[i] = letterBytes[rand.Intn(letterBytesLen)]
		}

	}
	return string(b)
}

func RandPassword(n int, no_special ...bool) string {
	//letterBytes := "abcdefhkmnprstuvwxyzABCEFGHIJKLMNPRSTUVWXYZ1234567890"
	letterUp := "ABCEFGHIJKLMNPRSTUVWXYZ"
	letterLower := "abcdefhkmnprstuvwxyz"
	letterDigits := "1234567890"
	letterSpecial := "!@#$%&*?"
	noSpecial := false
	if len(no_special) >= 1 {
		noSpecial = no_special[0]
	}

	password := RandString(1, letterUp) + RandString(1, letterLower) + RandString(1, letterDigits)
	lastN := n - 3
	if !noSpecial {
		password += RandString(1, letterSpecial)
		lastN -= 1
	}
	if lastN > 0 {
		password += RandString(lastN, letterUp, letterLower, letterDigits)
	}
	passwordBytes := []byte(password)
	for i := len(passwordBytes) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		passwordBytes[i], passwordBytes[num] = passwordBytes[num], passwordBytes[i]
	}

	return string(passwordBytes)
}
