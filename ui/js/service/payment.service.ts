// rxjs
import { Observable } from "rxjs";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class PaymentService extends HTTPService {
  private readonly paymentsUrlSegment: string = ENV.api + "/payments";

  getPaymentCharts(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.paymentsUrlSegment + "/charts");
  }

  refreshPayments(): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.paymentsUrlSegment + "/refresh");
  }
}
