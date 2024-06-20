package util

import (
	"encoding/base64"
	"strings"
)

// DecodeBase64String ...
func DecodeBase64String(str string) ([]byte, error) {
	var isRaw = !strings.HasSuffix(str, "=")
	if strings.Contains(str, "+") {
		if isRaw {
			return base64.RawStdEncoding.DecodeString(str)
		}
		return base64.StdEncoding.DecodeString(str)
	}
	if isRaw {
		return base64.RawURLEncoding.DecodeString(str)
	}
	return base64.URLEncoding.DecodeString(str)
}
