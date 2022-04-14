// lit element
import { customElement } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// google
import '@google-web-components/google-chart';
import { GoogleChart } from '@google-web-components/google-chart';

// memberdashboard
import {
  chartBlue,
  chartDarkGray,
  chartDarkGreen,
  chartRed,
} from '../../../shared/styles/colors';
import { showComponent } from '../../../shared/functions';
import { reportsChartStyle } from './report-chart.style';
import { ReportService } from '../../services/report.service';
import {
  ChartOptions,
  ReportChartResponse,
} from '../../types/api/reports-chart-response';
import { Inject } from '../../../shared/di/inject';

@customElement('report-chart')
export class ReportChart extends LitElement {
  @Inject('report')
  private reportService: ReportService;

  reportsCharts: Array<ReportChartResponse> = [];

  membershipTrendsData: ReportChartResponse;

  membershipDistributionData: Array<ReportChartResponse> = [];

  membershipTrendsChart: GoogleChart;

  membershipDistributionChart: GoogleChart;

  static get styles(): CSSResult[] {
    return [reportsChartStyle];
  }

  firstUpdated(): void {
    this.membershipTrendsChart =
      this.shadowRoot?.querySelector('#membership-trends');
    this.membershipDistributionChart = this.shadowRoot?.querySelector(
      '#membership-distribution'
    );
    this.createMemberReportChart();
  }

  creatMembershipTrendChart(chartData: ReportChartResponse): void {
    this.membershipTrendsChart.type = chartData.type;
    this.membershipTrendsChart.options = chartData.options;
    this.membershipTrendsChart.rows = chartData.rows;
    this.membershipTrendsChart.cols = chartData.cols;
  }

  createMembershipDistributionChart(chartData: ReportChartResponse): void {
    const options: ChartOptions = { ...chartData.options };
    options.title += ' - Membership Distribution';

    this.membershipDistributionChart.type = chartData.type;
    this.membershipDistributionChart.options = options;
    this.membershipDistributionChart.rows = chartData.rows;
    this.membershipDistributionChart.cols = chartData.cols;
  }

  openMembershipMonthsOptions(): void {
    this.requestUpdate();
    showComponent('#membership-months', this.shadowRoot);
  }

  createMemberReportChart(): void {
    this.reportService.getReportsCharts().subscribe({
      next: (result: ReportChartResponse[]) => {
        this.reportsCharts = result;
        this.reportsCharts = this.updateReportsCharts(this.reportsCharts);

        [this.membershipTrendsData] = this.reportsCharts;
        this.membershipDistributionData = this.reportsCharts.splice(1);

        this.creatMembershipTrendChart(this.membershipTrendsData);
        this.createMembershipDistributionChart(
          this.membershipDistributionData[
            this.membershipDistributionData.length - 1
          ]
        );
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to create report chart');
      },
    });
  }

  updateReportsCharts(
    reportsCharts: ReportChartResponse[]
  ): ReportChartResponse[] {
    reportsCharts.forEach((x: ReportChartResponse) => {
      x.options.colors = [
        chartBlue.cssText,
        chartDarkGreen.cssText,
        chartRed.cssText,
        chartDarkGray.cssText,
      ];
    });

    return reportsCharts;
  }

  updateMembershipDistributionChart(chartData: ReportChartResponse): void {
    this.createMembershipDistributionChart(chartData);
  }

  getReportMonthOptions(): TemplateResult {
    return html`
      <mwc-menu x="10" y="40" id="membership-months">
        ${this.membershipDistributionData.map((x: ReportChartResponse) => {
          return html`
            <mwc-list-item
              @click=${() => this.updateMembershipDistributionChart(x)}
            >
              ${x.options.title}
            </mwc-list-item>
          `;
        })}
      </mwc-menu>
    `;
  }

  render(): TemplateResult {
    return html`
      <div>
        <div class="select-month">
          <mwc-button
            raised
            label="Select a month"
            @click=${this.openMembershipMonthsOptions}
          ></mwc-button>
          ${this.getReportMonthOptions()}
        </div>
        <div id="report-chart-container">
          <google-chart id="membership-distribution"> </google-chart>
          <google-chart id="membership-trends"> </google-chart>
        </div>
      </div>
    `;
  }
}
