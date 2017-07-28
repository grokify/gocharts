package c3

type C3Chart struct {
	Bindto string      `json:"bindto,omitempty"`
	Data   C3ChartData `json:"data,omitempty"`
	Donut  C3Donut     `json:"donut,omitempty"`
}

type C3ChartData struct {
	Columns [][]interface{} `json:"columns,omitempty"`
	Type    string          `json:"type,omitempty"`
}

type C3Donut struct {
	Title string `json:"title,omitempty"`
}
