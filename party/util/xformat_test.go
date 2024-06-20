package util

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateRandNByLen(t *testing.T) {
	convey.Convey("Normal Result", t, func() {
		n := GenerateRandNByLen(0)
		convey.So(n, convey.ShouldEqual, 0)
		n = GenerateRandNByLen(1)
		convey.So(n, convey.ShouldEqual, 10)
		n = GenerateRandNByLen(2)
		convey.So(n, convey.ShouldEqual, 100)
		n = GenerateRandNByLen(4)
		convey.So(n, convey.ShouldEqual, 10000)
	})
}

func TestMustGetFreePort(t *testing.T) {
	convey.Convey("Normal Result", t, func() {
		port, err := GetFreePort()
		if err != nil {
			t.Fatal(err)
		}
		convey.So(port, convey.ShouldBeGreaterThan, 0)
	})
}

func TestVerifyMobileFormat(t *testing.T) {
	convey.Convey("Normal Result", t, func() {
		ok := VerifyMobileFormat("11111111111", MostLoose)
		convey.So(ok, convey.ShouldEqual, true)

		ok = VerifyMobileFormat("21111111111", MostLoose)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("13333333333", Loose)
		convey.So(ok, convey.ShouldEqual, true)

		ok = VerifyMobileFormat("19333333333", Loose)
		convey.So(ok, convey.ShouldEqual, true)

		ok = VerifyMobileFormat("12333333333", Loose)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("11333333333", Loose)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("13333333333", Rigor)
		convey.So(ok, convey.ShouldEqual, true)

		ok = VerifyMobileFormat("19533333333", Rigor)
		convey.So(ok, convey.ShouldEqual, true)

		ok = VerifyMobileFormat("16233333333", Rigor)
		convey.So(ok, convey.ShouldEqual, true)

		ok = VerifyMobileFormat("16333333333", Rigor)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("16933333333", Rigor)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("15433333333", Rigor)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("16433333333", Rigor)
		convey.So(ok, convey.ShouldEqual, false)

		ok = VerifyMobileFormat("19433333333", Rigor)
		convey.So(ok, convey.ShouldEqual, false)
	})
}

func TestWeightedChoice(t *testing.T) {
	convey.Convey("Normal Result", t, func() {
		choices := []*Choice{
			{Weight: 5, Item: "5%权重"},
			{Weight: 60, Item: "60%权重"},
			{Weight: 20, Item: "20%权重"},
			{Weight: 15, Item: "15%权重"},
		}
		var (
			choicesed = make([]*Choice, 0)
			cnt       = 1000
		)
		for i := 0; i < cnt; i++ {
			c, err := WeightedRandom(choices)
			if err != nil {
				t.Fatal(err)
			}
			choicesed = append(choicesed, c)
		}

		var sumWeight uint32
		for _, choice := range choices {
			sumWeight += choice.Weight
		}

		for _, choice := range choices {
			var count int
			for _, cd := range choicesed {
				if choice.Item == cd.Item {
					count++
				}
			}
			t.Logf("%s次数: %d，占比: %.2f%s", choice.Item, count, float64(count)/float64(cnt)*100, "%")
		}
	})
}
