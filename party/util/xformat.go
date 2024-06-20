package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastrand"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	// MaxUint 无符号整型 uint 最大值
	MaxUint uint = ^uint(0)
	// MinUint 无符号整型 uint 最小值
	MinUint uint = 0
	// MaxInt 有符号整型 int 最大值
	MaxInt int = int(^uint(0) >> 1)
	// MinInt 有符号整型 int 最小值
	MinInt int = ^MaxInt
	// MaxInt32 有符号整型 int32 最大值
	MaxInt32 int32 = int32(^uint32(0) >> 1)
	// MinInt32 有符号整型 int32 最小值
	MinInt32 int32 = ^MaxInt32
)

// GenerateRandNByLen ...
func GenerateRandNByLen(len int) int {
	if len == 0 {
		return 0
	}
	ns := "1" + fmt.Sprintf("%0"+strconv.Itoa(len)+"d", 0)
	return StringToInt(ns)
}

// FormatFloat format float64
// The precision prec controls the number of digits
// The special precision -1 uses the smallest number of digits
func FormatFloat(num float64, prec int) float64 {
	s := strconv.FormatFloat(num, 'f', prec, 64)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return float64(-1)
	}
	return f
}

// Float64ToInt float64 to int
func Float64ToInt(num float64) int {
	s := strconv.FormatFloat(num, 'f', 0, 64)
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

// CreatePathIfNotExists create path if it was not existed
func CreatePathIfNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

// IntToBool int to type bool
func IntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}

// RandInt64 the scope of rand from min to max
func RandInt64(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}

// Response ...
// NOTE: *grequests.Response 内部将数值型处理成 float64, 所以此处返回的数值型基本为 float64
func Response(response *grequests.Response, errMsg string) (m map[string]interface{}, err error) {
	if response.StatusCode != 200 || !response.Ok {
		return nil, errors.New(errMsg)
	}
	err = json.Unmarshal(response.Bytes(), &m)
	if err != nil {
		err = fmt.Errorf("response: %s, error: %w", string(response.Bytes()), err)
	}
	return
}

// ResponseMap type grequests.Response of http response to type map[string]interface{}
// NOTE: 该方法因为 return simplejson.Json.Map(), 所以所有数值型都是 json.Number 类型
func ResponseMap(response *grequests.Response, errMsg string) (map[string]interface{}, error) {
	if response.StatusCode != 200 || !response.Ok {
		return nil, errors.New(errMsg)
	}
	// NOTE： response is closed in response.Bytes() [(*Response).populateResponseByteBuffer]
	_json, err := simplejson.NewJson(response.Bytes())
	if err != nil {
		return nil, err
	}
	return _json.Map()
}

// SimplejsonToBytes simplejson to type []byte
func SimplejsonToBytes(j *simplejson.Json) ([]byte, error) {
	if j == nil {
		return nil, errors.New("param j is nil")
	}
	return json.Marshal(j.MustMap())
}

// MD5WithSalt ...
func MD5WithSalt(toSign, salt string) string {
	h := md5.New()
	h.Write([]byte(toSign))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5 ...
func MD5(toSign string) string {
	return MD5WithSalt(toSign, "")
}

// Sha1WithSalt ...
func Sha1WithSalt(data, salt string) []byte {
	h := sha1.New()
	h.Write([]byte(data))
	h.Write([]byte(salt))
	return h.Sum(nil)
}

// Sha1 ...
func Sha1(data string) []byte {
	return Sha1WithSalt(data, "")
}

// Sha256WithSalt ...
func Sha256WithSalt(data, salt string) []byte {
	h := sha256.New()
	h.Write([]byte(data))
	h.Write([]byte(salt))
	return h.Sum(nil)
}

// Sha256 ...
func Sha256(data string) []byte {
	return Sha256WithSalt(data, "")
}

// FourAfterMobile 海外地区手机号规则: 区号-手机号, 如: 626-3210900 (美国)，中国内地省去区号, 如: 13800138000
func FourAfterMobile(mobile string) string {
	ms := strings.Split(mobile, "-")
	if mobile = ms[len(ms)-1]; len(mobile) < 4 {
		return ""
	}
	return mobile[len(mobile)-4:]
}

// UnicodeEmojiDecode 表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

// UnicodeEmojiCode 表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

// FilterEmoji ...
func FilterEmoji(content string) string {
	newContent := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newContent += string(value)
		}
	}
	return newContent
}

// GbkToUtf8 ...
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk ...
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// FuzzyMobile for example, change 13800138000 to 138****8000
func FuzzyMobile(mobile string) string {
	if len(mobile) < 8 {
		return mobile
	}
	fuzzyNum := mobile[3:7]
	return strings.Replace(mobile, fuzzyNum, "****", -1)
}

// MapToQueryString returns unescape querystring
// map[string]interface{}{"bar": "baz", "foo": "quux"} -> "bar=baz&foo=quux"
func MapToQueryString(m map[string]interface{}) string {
	uv := url.Values{}
	for k, v := range m {
		uv.Add(k, fmt.Sprint(v))
	}
	return uv.Encode()
}

// MapToEscapeQueryString returns escape querystring
// map[string]interface{}{"bar": "baz", "foo": "quux"} -> "bar%3Dbaz%26foo%3Dquux"
func MapToEscapeQueryString(m map[string]interface{}) string {
	return url.QueryEscape(MapToQueryString(m))
}

// QueryStringToMap unescape querystring to map
// "bar=baz&foo=quux" -> map[string]interface{}{"bar": "baz", "foo": "quux"}
func QueryStringToMap(es string) (m map[string]interface{}, err error) {
	uv := url.Values{}
	if uv, err = url.ParseQuery(es); err != nil {
		return
	}
	m = make(map[string]interface{}, len(uv))
	for k, v := range uv {
		m[k] = strings.Join(v, ",")
	}
	return
}

// EscapeQueryStringToMap escape querystring to map
// "bar%3Dbaz%26foo%3Dquux" -> map[string]interface{}{"bar": "baz", "foo": "quux"}
func EscapeQueryStringToMap(ues string) (m map[string]interface{}, err error) {
	var es string
	if es, err = url.QueryUnescape(ues); err != nil {
		return
	}
	return QueryStringToMap(es)
}

// MustMarshal ...
func MustMarshal(v interface{}, isEscape ...bool) (data []byte) {
	buf, _ := marshal(v, isEscape...)
	data = buf.Bytes()
	return
}

// MustMarshalToString ...
func MustMarshalToString(v interface{}, isEscape ...bool) (str string) {
	buf, _ := marshal(v, isEscape...)
	str = buf.String()
	return
}

// Marshal ...
func Marshal(v interface{}, isEscape ...bool) (data []byte, err error) {
	buf, err := marshal(v, isEscape...)
	if err != nil {
		return
	}
	data = buf.Bytes()
	return
}

// marshal ...
func marshal(v interface{}, isEscape ...bool) (buf *bytes.Buffer, err error) {
	buf = bytes.NewBuffer([]byte{})
	jsonEncoder := jsoniter.NewEncoder(buf)
	if len(isEscape) > 0 {
		jsonEncoder.SetEscapeHTML(isEscape[0])
	}
	err = jsonEncoder.Encode(v)
	// NOTE: fix bug: jsonEncoder.Encode 默认在最后加了个 \n
	buf = bytes.NewBuffer(bytes.TrimSuffix(buf.Bytes(), []byte("\n")))
	return
}

type VerifyMobileLevel uint32

const (
	Loose     VerifyMobileLevel = iota // 宽松， 只要 13、14、15、16、17、18、19 开头即可
	Rigor                              // 严谨， 根据工信部公布的手机号段校验
	MostLoose                          // 最宽松，只要1开头即可
)

// VerifyMobileFormat ...
func VerifyMobileFormat(mobileNum string, level VerifyMobileLevel) bool {
	var regular string
	switch level {
	case Loose:
		regular = "^13|14|15|16|17|18|19\\d{9}$"
	case Rigor:
		// NOTE: 2022年4月版 https://www.qqzeng-ip.com/res/phone-hf.png
		regular = "^((13[0-9])|(14[0-9])|(15[0-3,5-9])|(16[2,5-7])|(17[0-9])|(18[0-9])|(19[0-3,5-9]))\\d{8}$"
	case MostLoose:
		regular = "^1\\d{10}$"
	default:
		return false
	}
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

type (
	Choices []*Choice
	Choice  struct {
		Weight uint32
		Item   interface{}
	}
)

// Len Less Swap ...
func (c Choices) Len() int           { return len(c) }
func (c Choices) Less(i, j int) bool { return c[i].Weight < c[j].Weight }
func (c Choices) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// WeightedRandom ...
func WeightedRandom(choices Choices) (choice *Choice, err error) {
	var sumWeight uint32
	for _, c := range choices {
		sumWeight += c.Weight
	}
	if sumWeight == 0 {
		err = errors.New("sum weight cannot be 0")
		return
	}

	// NOTE: rand.Int31n() 存在性能问题，并发量大的情况下尽量不要使用
	result := fastrand.Uint32n(sumWeight)
	sort.Sort(choices)
	for _, _choice := range choices {
		// NOTE: 注意 uint32 减法
		if result < _choice.Weight {
			choice = _choice
			break
		}
		if result >= _choice.Weight {
			result -= _choice.Weight
		}
	}
	if choice == nil {
		err = errors.New("unknown error: weighted random choice is nil")
		return
	}
	return
}
