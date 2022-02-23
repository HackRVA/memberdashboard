export interface ChartOptions {
  title: string;
  pieHole: number;
  legend: string;
  curveType: string;
  colors: string[];
}

export interface ChartCols {
  label: string;
  type: string;
}

export interface ReportChartResponse {
  id: string;
  type: string;
  options: ChartOptions;
  rows: any[]; // returns an array of arrays of string and number
  cols: ChartCols[];
}
