// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { ENV } from '../../env';
import { Injectable } from '../../shared/di/injectable';
import { HTTPService } from '../../shared/services/http.service';
import { ReportChartResponse } from '../types/api/reports-chart-response';
import { ChurnResponse } from '../types/api/churn-response';

@Injectable('report')
export class ReportService extends HTTPService {
  private readonly reportsUrlSegment: string = ENV.api + '/reports';

  getMemberChurn(): Observable<ChurnResponse> {
    return this.get<ChurnResponse>(this.reportsUrlSegment + '/churn');
  }

  getReportsCharts(): Observable<ReportChartResponse[]> {
    return this.get<ReportChartResponse[]>(
      this.reportsUrlSegment + '/membercounts'
    );
  }

  getAccessChart(resourceName: string): Observable<ReportChartResponse> {
    return this.get<ReportChartResponse>(
      `${this.reportsUrlSegment}/access?resourceName=${resourceName}`
    );
  }
}
