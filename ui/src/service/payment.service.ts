// rxjs
import { Observable } from "rxjs";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";
import { PaymentChartResponse } from "../components/payments/types";

export class PaymentService extends HTTPService {
  private readonly paymentsUrlSegment: string = ENV.api + "/payments";

  getPaymentCharts(): Observable<PaymentChartResponse[]> {
    return this.get<PaymentChartResponse[]>(
      this.paymentsUrlSegment + "/charts"
    );
  }

  refreshPayments(): Observable<void> {
    return this.post<void>(this.paymentsUrlSegment + "/refresh");
  }
}
