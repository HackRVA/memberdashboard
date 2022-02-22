// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { ENV } from '../../env';
import { Injectable } from '../../shared/di/injectable';
import { HTTPService } from '../../shared/services/http.service';
import { PaymentChartResponse } from '../types/api/payment-chart-response';

@Injectable('payment')
export class PaymentService extends HTTPService {
  private readonly paymentsUrlSegment: string = ENV.api + '/reports';

  getPaymentCharts(): Observable<PaymentChartResponse[]> {
    return this.get<PaymentChartResponse[]>(
      this.paymentsUrlSegment + '/membercounts'
    );
  }
}
