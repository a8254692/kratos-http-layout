package util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//Uniq delete duplicate data in the array of ls
func Uniq(ls []int) []int {
	intInSlice := func(i int, list []int) bool {
		for _, v := range list {
			if v == i {
				return true
			}
		}
		return false
	}
	var Uniq []int
	for _, v := range ls {
		if !intInSlice(v, Uniq) {
			Uniq = append(Uniq, v)
		}
	}
	return Uniq
}

// FindIndex find value in the array of ls.
// Returns the position of value in ls if it has been found.
// Returns -1 if it hasn's been found.
func FindIndex(ls []int, value int) int {
	for k, v := range ls {
		if value == v {
			return k
		}
	}
	return -1
}

// IntArrayToString int arrays to string replaced by spe
func IntArrayToString(arrays []int, spe string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(arrays), " ", spe, -1), "[]")
}

// Contains reports whether target is within obj
// target supports arrary,slice,map
// but performance is not good
func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}

// ContainsInt ...
func ContainsInt(obj int, target []int) int {
	for _, v := range target {
		if v == obj {
			return obj
		}
	}
	return MinInt
}

// ContainsInt32 ...
func ContainsInt32(obj int32, target []int32) int32 {
	for _, v := range target {
		if v == obj {
			return obj
		}
	}
	return int32(MinInt32)
}

// ContainsString ...
func ContainsString(obj string, target []string) string {
	for _, v := range target {
		if v == obj {
			return obj
		}
	}
	return ""
}

// IsContainsString ...
func IsContainsString(obj string, target []string) bool {
	for _, v := range target {
		if v == obj {
			return true
		}
	}
	return false
}

// StringSliceRemove string slice remove s element
func StringSliceRemove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// StringSliceRemoveElement ...
func StringSliceRemoveElement(slice []string, s string) []string {
	result := make([]string, 0, len(slice))
	for _, s2 := range slice {
		if s2 != s {
			result = append(result, s2)
		}
	}
	return result
}

// Int64SliceRemoveElement ...
func Int64SliceRemoveElement(slice []int64, s int64) []int64 {
	result := make([]int64, 0, len(slice))
	for _, s2 := range slice {
		if s2 != s {
			result = append(result, s2)
		}
	}
	return result
}

// StringSliceRemoveRep ...
func StringSliceRemoveRep(slice []string) []string {
	result := make([]string, 0, len(slice))
	temp := map[string]struct{}{}
	for _, item := range slice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
