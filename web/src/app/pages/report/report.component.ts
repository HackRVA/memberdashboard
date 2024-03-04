import '@google-web-components/google-chart';
import {
  Component,
  DestroyRef,
  ElementRef,
  OnInit,
  Renderer2,
  ViewChild,
  inject,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { Observable, forkJoin, of, switchMap, tap } from 'rxjs';
import { ReportService } from './services';
import { ChartOptions, ReportChartResponse } from './types';

@Component({
  selector: 'md-report',
  standalone: true,
  imports: [MatButtonModule, MatMenuModule],
  providers: [ReportService],
  templateUrl: './report.component.html',
  styleUrl: './report.component.scss',
})
export class ReportComponent implements OnInit {
  @ViewChild('distribution') distributionChart: ElementRef;
  @ViewChild('trend') trendChart: ElementRef;
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  memberChurn: number = 0;
  membershipDistributionData: ReportChartResponse[] = [];

  constructor(
    private readonly reportService: ReportService,
    private readonly renderer: Renderer2
  ) {}

  ngOnInit(): void {
    this.fetchAndLoadReport().subscribe();
  }

  private fetchAndLoadReport(): Observable<void> {
    return forkJoin({
      churnResponse: this.reportService.getMemberChurn(),
      reportsResponse: this.reportService.getReportsCharts(),
    }).pipe(
      takeUntilDestroyed(this._destroyRef),
      tap(({ churnResponse, reportsResponse }) => {
        if (churnResponse) {
          this.memberChurn = churnResponse.churn;
        }

        if (reportsResponse) {
          this.membershipDistributionData = reportsResponse.splice(1); // ignoring the first element since it's only for trend data
          this.createMembershipDistributionChart(
            this.membershipDistributionData[
              this.membershipDistributionData.length - 1
            ]
          );
          this.createMembershipTrendChart(reportsResponse[0]);
        }
      }),
      switchMap(() => of(null))
    );
  }

  updateMembershipDistributionChart(chartData: ReportChartResponse): void {
    const existingGoogleChart =
      this.distributionChart.nativeElement.querySelector('google-chart');
    if (existingGoogleChart) {
      this.renderer.removeChild(
        this.distributionChart.nativeElement,
        existingGoogleChart
      );
    }
    this.createMembershipDistributionChart(chartData);
  }

  private createMembershipDistributionChart(
    chartData: ReportChartResponse
  ): void {
    const options: ChartOptions = { ...chartData.options };
    options.title += ' - Membership Distribution';
    options.colors = ['#6200ee', '#e9437a', '#888888', '#50c878'];

    this.createGoogleChart(
      this.distributionChart,
      Object.assign({}, chartData, { options: options })
    );
  }

  private createMembershipTrendChart(chartData: ReportChartResponse): void {
    const options: ChartOptions = { ...chartData.options };
    options.colors = ['#6200ee'];

    this.createGoogleChart(
      this.trendChart,
      Object.assign({}, chartData, { options: options })
    );
  }

  private createGoogleChart(
    el: ElementRef,
    chartData: ReportChartResponse
  ): void {
    const googleChart = this.renderer.createElement('google-chart');
    this.renderer.setAttribute(
      googleChart,
      'options',
      JSON.stringify(chartData.options)
    );
    this.renderer.setAttribute(googleChart, 'type', chartData.type);
    this.renderer.setAttribute(
      googleChart,
      'rows',
      JSON.stringify(chartData.rows)
    );
    this.renderer.setAttribute(
      googleChart,
      'cols',
      JSON.stringify(chartData.cols)
    );
    this.renderer.appendChild(el.nativeElement, googleChart);
  }
}
