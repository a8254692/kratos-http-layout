package util

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// StringToInt string to int. Returns -1 If it has error
func StringToInt(value string) int {
	v, e := strconv.Atoi(value)
	if e != nil {
		return -1
	}
	return v
}

// StringToMap string to map[string]interface{}. Returns empty map If it has error
func StringToMap(value string) (m map[string]interface{}) {
	m = make(map[string]interface{})
	_ = json.Unmarshal([]byte(value), &m)
	return
}

// StringToFloat64 string to float64. Returns -1 If it has error
func StringToFloat64(value string) float64 {
	v, e := strconv.ParseFloat(value, 10)
	if e != nil {
		return -1
	}
	return v
}

// Substr ...
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	if start < 0 || end < 0 || start > len(rs) || end > len(rs) || start > end {
		return ""
	}
	return string(rs[start:end])
}

// Float64ToString float64 to string
// The precision prec controls the number of digits
// The special precision -1 uses the smallest number of digits
func Float64ToString(num float64, prec int) string {
	return strconv.FormatFloat(num, 'f', prec, 64)
}

// Int64ToString int64 into string
func Int64ToString(num int64) string {
	//10 为十进制
	return strconv.FormatInt(num, 10)
}

// StringBuilder ...
func StringBuilder(s ...string) string {
	b := strings.Builder{}
	b.Grow(128) // pre-allocated 128 bytes
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

// StringToByte ...
func StringToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// BytesToString ...
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// CamelToUnderline ...
func CamelToUnderline(s string) string {
	buffer := new(bytes.Buffer)
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

type RandStrType string

// Append ...
func (r RandStrType) Append(rt RandStrType) RandStrType {
	return r + rt
}

const (
	RandStrNumber RandStrType = "0123456789"
	RandStrLower  RandStrType = "abcdefghijklmnopqrstuvwxyz"
	RandStrUpper  RandStrType = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandStr ...
type RandStr struct {
	letters      RandStrType
	letterIdBits int
	letterIdMask int64
	letterIdMax  int
}

// NewRandStr ...
func NewRandStr(rt RandStrType) *RandStr {
	// 6 bits to represent a letter index
	var letterIdBits = 6
	return &RandStr{
		letters:      rt,
		letterIdBits: letterIdBits,
		// All 1-bits as many as letterIdBits
		letterIdMask: 1<<letterIdBits - 1,
		letterIdMax:  63 / letterIdBits,
	}
}

// Rand copy form https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func (r *RandStr) Rand(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, randSrc.Int63(), r.letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), r.letterIdMax
		}
		if idx := int(cache & r.letterIdMask); idx < len(r.letters) {
			b[i] = r.letters[idx]
			i--
		}
		cache >>= r.letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// EscapeSql ...
func EscapeSql(str string) string {
	return strings.NewReplacer(`\`, `\\`, `_`, `\_`, `%`, `\%`, `'`, `\'`).Replace(str)
}
