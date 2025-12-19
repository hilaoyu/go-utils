package utilRandom

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RandReader() io.Reader {
	return random
}

func RandInt64(n int64) int64 {
	return random.Int63n(n)
}
func RandInt(n int) int {
	return random.Intn(n)
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

func UniqId(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	return fmt.Sprintf(
		"%s%x-%x",
		prefix,
		time.Now().UnixNano(),
		rand.Uint32(),
	)
}
