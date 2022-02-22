// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { ENV } from '../../env';
import { Injectable } from '../../shared/di/injectable';
import { HTTPService } from '../../shared/services/http.service';
import { ReportChartResponse } from '../types/api/reports-chart-response';

@Injectable('report')
export class ReportService extends HTTPService {
  private readonly reportsUrlSegment: string = ENV.api + '/reports';

  getReportsCharts(): Observable<ReportChartResponse[]> {
    return this.get<ReportChartResponse[]>(
      this.reportsUrlSegment + '/membercounts'
    );
  }
}
