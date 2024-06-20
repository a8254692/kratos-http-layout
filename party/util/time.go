package util

import (
	"time"
)

// TimeFormat format the value of type time.Time from inLayout to outLayout
// value: the value of time to format
// inLayout: input layout of the value
// outLayout: output layout of the value
func TimeFormat(value, inLayout, outLayout string) (valuestr string, err error) {
	var _value time.Time
	if _value, err = time.Parse(inLayout, value); err != nil {
		return
	}
	valuestr = _value.Format(outLayout)
	return
}

// LocalTimeFormat ...
func LocalTimeFormat(value, inLayout, outLayout string) (valuestr string, err error) {
	var _value time.Time
	if _value, err = time.ParseInLocation(inLayout, value, time.Local); err != nil {
		return
	}
	valuestr = _value.Format(outLayout)
	return
}

// TimeStringToUnix ...
func TimeStringToUnix(value, inLayout string) int {
	_time, err := time.Parse(inLayout, value)
	if err != nil {
		return 0
	}
	return int(_time.Unix())
}

// LocalTimeStringToUnix ...
func LocalTimeStringToUnix(value, inLayout string) int {
	_time, err := time.ParseInLocation(inLayout, value, time.Local)
	if err != nil {
		return 0
	}
	return int(_time.Unix())
}

// TimeRankScore ...
func TimeRankScore() float64 {
	return float64(time.Now().Local().Unix())
}

// GetZeroTime ...
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetFirstDateOfMonth ...
func GetFirstDateOfMonth(d time.Time) time.Time {
	return GetZeroTime(d.AddDate(0, 0, -d.Day()+1))
}

// GetLastDateOfMonth ...
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// UnixTimeToRFC3339 ...
func UnixTimeToRFC3339(sec int64, nsec int64) string {
	if sec == 0 {
		return ""
	}
	return time.Unix(sec, nsec).Format(time.RFC3339)
}

// TimeRFC3339ToUnix ...
func TimeRFC3339ToUnix(format string) (int64, error) {
	t, err := time.Parse(time.RFC3339, format)
	return t.Unix(), err
}

// SimpleTimeToUnix ...
func SimpleTimeToUnix(format string) (int64, error) {
	t, err := time.ParseInLocation(time.DateTime, format, time.Local)
	return t.Unix(), err
}

// ParseTime ...
func ParseTime(value string, layouts ...string) (time.Time, error) {
	if len(layouts) == 0 {
		layouts = []string{time.DateTime, time.RFC3339}
	}
	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}
