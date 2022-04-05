package slot

import (
	"testing"
)

var bucketTests = []struct {
	size  int64
	val   int64
	index int64
	min   int64
	max   int64
}{
	{5, 1, 0, 1, 5},
	{5, 2, 0, 1, 5},
	{5, 3, 0, 1, 5},
	{5, 4, 0, 1, 5},
	{5, 5, 0, 1, 5},
	{5, 6, 1, 6, 10},
	{5, 7, 1, 6, 10},
	{5, 8, 1, 6, 10},
	{5, 9, 1, 6, 10},
	{5, 10, 1, 6, 10},
	{5, 11, 2, 11, 15}}

func TestBucket(t *testing.T) {
	for _, tt := range bucketTests {
		min, max := BucketMinMax(tt.size, tt.index)
		if min != tt.min || max != tt.max {
			t.Errorf("data.BucketMinMax: with [%v,%v] want [%v,%v] got [%v,%v]",
				tt.size, tt.index, tt.min, tt.max, min, max)
		}

		index := BucketIndex(tt.size, tt.val)
		if index != tt.index {
			t.Errorf("data.BucketIndex: with [%v,%v] want [%v] got [%v]",
				tt.size, tt.val, tt.index, index)
		}
	}
}
