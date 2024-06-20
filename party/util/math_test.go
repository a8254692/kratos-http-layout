package util

import (
	"testing"
)

func TestRound(t *testing.T) {
	testCases := []struct {
		name string
		req  float64
		n    int32
		res  float64
		want bool
	}{
		{
			name: "t0",
			req:  1.2251233467,
			n:    2,
			res:  1.23,
			want: true,
		},
		{
			name: "t1",
			req:  0.825,
			n:    2,
			res:  0.83,
			want: true,
		},
		{
			name: "t2",
			req:  0.825,
			n:    1,
			res:  0.8,
			want: true,
		},
		{
			name: "t3",
			req:  0.825,
			n:    0,
			res:  1,
			want: true,
		},
		{
			name: "t4",
			req:  0.84,
			n:    2,
			res:  0.84,
			want: true,
		},
		{
			name: "t5",
			req:  1.23,
			n:    1,
			res:  1.23,
			want: false,
		},
		{
			name: "t6",
			req:  1.224,
			n:    2,
			res:  1.22,
			want: true,
		},
		{
			name: "t7",
			req:  1.229,
			n:    2,
			res:  1.23,
			want: true,
		},
		{
			name: "t8",
			req:  1.220,
			n:    2,
			res:  1.22,
			want: true,
		},
		{
			name: "t9",
			req:  1.0,
			n:    0,
			res:  1,
			want: true,
		},
		{
			name: "t10",
			req:  1.345,
			n:    0,
			res:  1,
			want: true,
		},
		{
			name: "t11",
			req:  0.145,
			n:    2,
			res:  0.15,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rounding := Round(tc.req, tc.n)
			if tc.want {
				if rounding != tc.res {
					t.Errorf("req %v, n %d, got %v, expected %v, want %v", tc.req, tc.n, rounding, tc.res, tc.want)
					return
				}
			} else {
				if rounding == tc.res {
					t.Errorf("req %v, n %d, got %v, expected %v, want %v", tc.req, tc.n, rounding, tc.res, tc.want)
					return
				}
			}
			t.Logf("req %v, n %d, got %v, expected %v, want %v", tc.req, tc.n, rounding, tc.res, tc.want)
		})
	}
}
