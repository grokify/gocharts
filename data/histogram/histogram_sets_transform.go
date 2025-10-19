package histogram

var funcStringTransformNoop = func(s string) string { return s }

func (hsets *HistogramSets) FilterHsetNames(inFunc, exFunc func(name string) bool) *HistogramSets {
	out := NewHistogramSets(hsets.Name)
	for hsetName, hset := range hsets.Items {
		if (exFunc != nil && exFunc(hsetName)) ||
			inFunc == nil ||
			!inFunc(hsetName) {
			continue
		}
		for histName, hist := range hset.Items {
			for binName, binSum := range hist.Items {
				out.Add(
					hsetName,
					histName,
					binName,
					binSum,
					false,
				)
			}
		}
	}
	return out
}

func (hsets *HistogramSets) TransformBinNamesMap(xfBinNames map[string]string) *HistogramSets {
	return hsets.TransformNames(
		nil,
		nil,
		func(name string) string {
			if newName, ok := xfBinNames[name]; ok {
				return newName
			} else {
				return name
			}
		},
	)
}

func (hsets *HistogramSets) TransformNames(xfHsetName, xfFuncHistName, xfFuncBinName func(name string) string) *HistogramSets {
	out := NewHistogramSets(hsets.Name)
	if xfHsetName == nil {
		xfHsetName = funcStringTransformNoop
	}
	if xfFuncHistName == nil {
		xfFuncHistName = funcStringTransformNoop
	}
	if xfFuncBinName == nil {
		xfFuncBinName = funcStringTransformNoop
	}
	for hsetName, hset := range hsets.Items {
		for histName, hist := range hset.Items {
			for binName, binSum := range hist.Items {
				out.Add(
					xfHsetName(hsetName),
					xfFuncHistName(histName),
					xfFuncBinName(binName),
					binSum,
					false,
				)
			}
		}
	}
	return out
}
