package models

// ChartOptions -- config option for the chart
type ChartOptions struct {
	Title     string  `json:"title"`
	CurveType string  `json:"curveType"`
	PieHole   float64 `json:"pieHole"`
	Legend    string  `json:"legend"`
}

// ChartRow -- chart row info
type ChartRow struct {
	MemberLevel, MemberCount interface{}
}

// ChartCol -- chart col info
type ChartCol struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// PaymentChart - deliver information to the frontend so that
//   we can display a monthly payment chart
type PaymentChart struct {
	Options ChartOptions    `json:"options"`
	Type    string          `json:"type"`
	Rows    [][]interface{} `json:"rows"`
	Cols    []ChartCol      `json:"cols"`
}
