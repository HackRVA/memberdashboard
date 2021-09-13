// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { ENV } from '../../env';
import { HTTPService } from '../../shared/services/http.service';
import { PaymentChartResponse } from '../types/api/payment-chart-response';

export class PaymentService extends HTTPService {
  private readonly paymentsUrlSegment: string = ENV.api + '/payments';

  getPaymentCharts(): Observable<PaymentChartResponse[]> {
    return this.get<PaymentChartResponse[]>(
      this.paymentsUrlSegment + '/charts'
    );
  }

  refreshPayments(): Observable<void> {
    return this.post<void>(this.paymentsUrlSegment + '/refresh');
  }
}
