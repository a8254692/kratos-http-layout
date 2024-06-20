package util

import (
	"fmt"
	"testing"
	"time"
)

func TestGetLastDateOfMonth(t *testing.T) {
	tn := time.Now()
	fmt.Println(GetZeroTime(tn).String())
	fmt.Println(GetFirstDateOfMonth(tn).String())
	fmt.Println(GetLastDateOfMonth(tn).String())
}

func TestParseTime(t *testing.T) {
	v := "2020-10-09T11:13:26+08:00"
	tm, err := ParseTime(v)
	if err != nil {
		t.Log(err)
	}
	t.Log(tm.Format(time.RFC3339))

	v1 := "2020-10-09 11:13:26"
	tm1, err := ParseTime(v1)
	if err != nil {
		t.Log(err)
	}
	t.Log(tm1.Format(time.RFC3339))

	v2 := "2020-10-09 11:13:"
	tm2, err := ParseTime(v2)
	if err != nil {
		t.Log(err)
	}
	t.Log(tm2.Format(time.RFC3339))
}
