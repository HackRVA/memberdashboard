export type ChurnResponse = {
  churn: number;
};

export type ChartOptions = {
  title: string;
  pieHole: number;
  legend: string;
  curveType: string;
  colors: string[];
};

export type ChartCols = {
  label: string;
  type: string;
};

export type ReportChartResponse = {
  id: string;
  type: string;
  options: ChartOptions;
  rows: any[]; // returns an array of arrays of string and number
  cols: ChartCols[];
};
