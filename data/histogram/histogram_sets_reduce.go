package histogram

// HistogramSetOneTwo returns a `HistogramSet` where the output histogram
// sums are the sums across source histograms.
func (hsets *HistogramSets) HistogramSetOneTwo() *HistogramSet {
	out := NewHistogramSet("")
	for hsetName, hset := range hsets.Items {
		for histName, hist := range hset.Items {
			sum := 0
			for _, binSum := range hist.Items {
				sum += binSum
			}
			out.Add(hsetName, histName, sum)
		}
	}
	return out
}

// HistogramSetOneThree returns a `HistogramSet` where the output histogram
// names are the source binNames.
func (hsets *HistogramSets) HistogramSetOneThree() *HistogramSet {
	out := NewHistogramSet("")
	for hsetName, hset := range hsets.Items {
		for _, hist := range hset.Items {
			for binName, binSum := range hist.Items {
				out.Add(hsetName, binName, binSum)
			}
		}
	}
	return out
}
