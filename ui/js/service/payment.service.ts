import { Observable } from "rxjs";
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class PaymentService extends HTTPService {
  private readonly api: string | undefined = ENV.api;

  getPaymentCharts(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/payments/charts");
  }
}

export namespace PaymentService {
  interface ChartOptions {}
  export interface PaymentChartResponse {
    id?: string;
    type?: string;
    options: ChartOptions;
    rows: string;
    cols: string;
  }
}
