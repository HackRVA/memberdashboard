package models

import "time"

type MemberCount struct {
	Month    time.Time `json:"month"`
	Classic  int       `json:"classic"`
	Standard int       `json:"standard"`
	Premium  int       `json:"premium"`
	Credited int       `json:"credited"`
}

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

// ReportChart - deliver information to the frontend so that
//   we can display a monthly payment chart
type ReportChart struct {
	Options ChartOptions    `json:"options"`
	Type    string          `json:"type"`
	Rows    [][]interface{} `json:"rows"`
	Cols    []ChartCol      `json:"cols"`
}

type AccessStats struct {
	Date         time.Time `json:"date"`
	AccessCount  int       `json:"accessCount"`
	ResourceName string    `json:"resourceName"`
}

type MemberChurn struct {
	Churn int `json:"churn"`
}
