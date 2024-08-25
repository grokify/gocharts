package google

import (
	"github.com/grokify/mogo/encoding/jsonutil"
)

type DataTable [][]any

func (dt DataTable) MustJSON() []byte {
	return jsonutil.MustMarshalOrDefault(dt, []byte(jsonutil.EmptyArray))
}
