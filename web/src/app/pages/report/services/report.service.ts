import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ChurnResponse, ReportChartResponse } from '../types/report.types';

@Injectable()
export class ReportService {
  private readonly _reportsUrlSegment: string = '/api/reports';

  constructor(private readonly http: HttpClient) {}

  getMemberChurn(): Observable<ChurnResponse> {
    return this.http.get<ChurnResponse>(this._reportsUrlSegment + '/churn');
  }

  getReportsCharts(): Observable<ReportChartResponse[]> {
    return this.http.get<ReportChartResponse[]>(
      this._reportsUrlSegment + '/membercounts'
    );
  }
}
