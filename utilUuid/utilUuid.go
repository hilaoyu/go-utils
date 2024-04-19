package utilUuid

import (
	"github.com/gofrs/uuid"
	"regexp"
	"strings"
)

func UuidGenerate() string {
	uid, err := uuid.NewV7()
	if nil != err {
		return ""
	}

	return strings.ReplaceAll(uid.String(), "-", "")
}

func IsUuid(obj interface{}) bool {
	var regUuid = regexp.MustCompile("^[0-9A-Za-z]{8}-?[0-9A-Za-z]{4}-?[0-9A-Za-z]{4}-?[0-9A-Za-z]{4}-?[0-9A-Za-z]{12}$")

	value, ok := obj.(string)
	if !ok {
		return false
	}

	ret := regUuid.MatchString(value)

	return ret
}
