package slot

func BucketMinMax(size, index int64) (int64, int64) {
	max := (index + 1) * size
	min := max - size + 1
	return min, max
}

func BucketIndex(size, this int64) int64 {
	if size == 0 {
		panic("Bucket Size Cannot be Zero")
	}
	index := (this - 1) / size
	return index
}

type Bucket struct {
	Min int64
	Max int64
}

func BucketsForMinMax(size, min, max int64) []Bucket {
	if min > max {
		tmp := min
		min = max
		max = tmp
	}
	minBucketIndex := BucketIndex(size, min)
	maxBucketIndex := BucketIndex(size, max)
	buckets := []Bucket{}
	for this := minBucketIndex; this <= maxBucketIndex; this++ {
		bucketMin, bucketMax := BucketMinMax(size, this)
		buckets = append(buckets, Bucket{Min: bucketMin, Max: bucketMax})
	}
	return buckets
}
