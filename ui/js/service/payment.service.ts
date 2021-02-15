// rxjs
import { Observable } from "rxjs";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class PaymentService extends HTTPService {
  private readonly api: string = ENV.api;

  getPaymentCharts(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/payments/charts");
  }
}
