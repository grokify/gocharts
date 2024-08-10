package google

import "encoding/json"

type DataTable [][]any

func (dt DataTable) MustJSON() []byte {
	if bytes, err := json.Marshal(dt); err != nil {
		return []byte("[]")
	} else {
		return bytes
	}
}
