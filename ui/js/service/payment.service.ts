import { Observable } from "rxjs";
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class PaymentService extends HTTPService {
  private readonly api: string = ENV.api;

  getPaymentCharts(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/payments/charts");
  }
}

export namespace PaymentService {
  export interface ChartOptions {
    title: string;
    pieHole: number;
    legend: string;
    curveType: string;
  }

  export interface ChartCols {
    label: string;
    type: string;
  }

  export interface PaymentChartResponse {
    id: string;
    type: string;
    options: ChartOptions;
    rows: any[]; // returns an array of arrays of string and number
    cols: ChartCols[];
  }

  export interface ChartAttributes {
    id: string;
    type: string;
    options: string;
    rows: string;
    cols: string;
  }
}
