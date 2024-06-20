package util

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTransStructToKeyMap(t *testing.T) {
	convey.Convey("Normal Result", t, func() {
		l := 10000000
		var ss = make([]string, 0, l)
		for i := 0; i < l; i++ {
			s := NewRandStr(RandStrLower.Append(RandStrNumber)).Rand(10)
			ss = append(ss, s)
		}
		rs := StringSliceRemoveRep(ss)
		t.Log(rs[:10])
		convey.So(len(rs), convey.ShouldEqual, len(ss))
	})

}
