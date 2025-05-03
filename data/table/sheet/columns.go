package sheet

import (
	"fmt"
	"slices"
)

func SliceReplaceValueAtLettersMap(s []string, m map[string]string) ([]string, error) {
	out := slices.Clone(s)
	var err error
	for k, v := range m {
		out, err = SliceReplaceValueAtLetters(s, k, v)
		if err != nil {
			return []string{}, err
		}
	}
	return out, nil
}

func SliceReplaceValueAtLetters(s []string, letters, v string) ([]string, error) {
	idx, err := ColLettersToIndex(letters)
	if err != nil {
		return []string{}, err
	}
	if int(idx) >= len(s) {
		return []string{}, fmt.Errorf("index out of bounds index letter (%s) index int (%d) len input (%d)",
			letters, idx, len(s))
	}
	out := slices.Clone(s)
	out[idx] = v
	return out, nil
}
