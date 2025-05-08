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
	num, err := ColLettersToNumber(letters)
	if err != nil {
		return []string{}, err
	}
	if int(num) > len(s) {
		return []string{}, fmt.Errorf("index out of bounds index letter (%s) index int (%d) len input (%d)",
			letters, num-1, len(s))
	}
	out := slices.Clone(s)
	out[num-1] = v
	return out, nil
}
