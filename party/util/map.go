package util

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"reflect"
)

// MapToSimplejson map[string]interface{} to simplejson.Json
func MapToSimplejson(m map[string]interface{}) (*simplejson.Json, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(bytes)
}

// MapToString interface{} to type string
// map[string]interface{}{"a": "aa","b": "bb"} -> "{"a":"aa","b":"bb"}"
func MapToString(m interface{}) string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// StructToMap struct to type map[string]interface{]}
func StructToMap(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	var _map map[string]interface{}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &_map)
	return _map, err
}

// FormToMap ...
func FormToMap(form map[string][]string) map[string]interface{} {
	m := make(map[string]interface{})
	for k, values := range form {
		m[k] = values[0]
	}
	return m
}

// TransStructToKeyMap ...
func TransStructToKeyMap(slice interface{}, key string, keyMap interface{}) {
	vs := reflect.ValueOf(slice)
	if vs.Kind() != reflect.Slice {
		panic("invalid error type")
	}

	vsMap := reflect.ValueOf(keyMap)
	for i := 0; i < vs.Len(); i++ {
		ptr, elem := vs.Index(i), vs.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		vsMap.SetMapIndex(elem.FieldByName(key), ptr)
	}
	return
}

// MergeMap ...
// NOTE: ms cover dest which has the same key
func MergeMap(dest map[string]interface{}, ms ...map[string]interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{})
	for k, v := range dest {
		m[k] = v
	}
	for _, _m := range ms {
		for k, v := range _m {
			m[k] = v
		}
	}
	return
}
