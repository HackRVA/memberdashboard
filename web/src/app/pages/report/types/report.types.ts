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
  type: string;
  options: ChartOptions;
  rows: [string, number][];
  cols: ChartCols[];
};
