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
import { ResourceService } from '../../../resource/services/resource.service';
import { ResourceResponse } from '../../../resource/types/api/resource-response';
import {
  ChartOptions,
  ReportChartResponse,
} from '../../types/api/reports-chart-response';
import { Inject } from '../../../shared/di/inject';

@customElement('report-chart')
export class ReportChart extends LitElement {
  @Inject('report')
  private reportService: ReportService;

  @Inject('resource')
  private resourceService: ResourceService;

  reportsCharts: Array<ReportChartResponse> = [];

  accessChart: ReportChartResponse;

  membershipTrendsData: ReportChartResponse;

  membershipAccessTrendsData: ReportChartResponse;

  membershipDistributionData: Array<ReportChartResponse> = [];

  membershipTrendsChart: GoogleChart;

  membershipAccessTrendsChart: GoogleChart;

  membershipDistributionChart: GoogleChart;

  resourceName: string;

  resources: Array<ResourceResponse>;

  static get styles(): CSSResult[] {
    return [reportsChartStyle];
  }

  firstUpdated(): void {
    this.membershipTrendsChart =
      this.shadowRoot?.querySelector('#membership-trends');
    this.membershipAccessTrendsChart = this.shadowRoot?.querySelector(
      '#membership-access-trends'
    );
    this.membershipDistributionChart = this.shadowRoot?.querySelector(
      '#membership-distribution'
    );

    this.getResources();
    this.createMemberReportChart();
    this.createAccessReportChart(this.resourceName);
  }

  creatMembershipTrendChart(chartData: ReportChartResponse): void {
    this.membershipTrendsChart.type = chartData.type;
    this.membershipTrendsChart.options = chartData.options;
    this.membershipTrendsChart.rows = chartData.rows;
    this.membershipTrendsChart.cols = chartData.cols;
  }

  createMembershipTrendChart(chartData: ReportChartResponse): void {
    this.membershipAccessTrendsChart.type = chartData.type;
    this.membershipAccessTrendsChart.options = chartData.options;
    this.membershipAccessTrendsChart.rows = chartData.rows;
    this.membershipAccessTrendsChart.cols = chartData.cols;
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

  openResourceOptions(): void {
    this.requestUpdate();
    showComponent('#resource-names', this.shadowRoot);
  }

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (response: ResourceResponse[]) => {
        this.resources = response;
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to get resources');
      },
    });
  }

  createAccessReportChart(name): void {
    this.reportService.getAccessChart(name).subscribe({
      next: (result: ReportChartResponse) => {
        this.accessChart = result;
        this.accessChart = this.updateAccessChart(this.accessChart);
        this.membershipAccessTrendsData = this.accessChart;
        this.createMembershipTrendChart(this.membershipAccessTrendsData);
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to create access report chart');
      },
    });
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

  updateAccessChart(accessChart: ReportChartResponse): ReportChartResponse {
    accessChart.options.colors = [
      chartBlue.cssText,
      chartDarkGreen.cssText,
      chartRed.cssText,
      chartDarkGray.cssText,
    ];

    return accessChart;
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

  updateAccessTrendsChart(name: string): void {
    this.createAccessReportChart(name);
  }

  getReportMonthOptions(): TemplateResult {
    return html`
      <mwc-menu x="10" y="40" id="membership-months">
        ${this.membershipDistributionData.map(
          (x: ReportChartResponse) => html`
            <mwc-list-item
              @click=${() => this.updateMembershipDistributionChart(x)}
            >
              ${x.options.title}
            </mwc-list-item>
          `
        )}
      </mwc-menu>
    `;
  }

  getResourceOptions(): TemplateResult {
    return html`
      <mwc-menu x="10" y="40" id="resource-names">
        ${this.resources?.map(
          (x: ResourceResponse) => html`
            <mwc-list-item @click=${() => this.updateAccessTrendsChart(x.name)}>
              ${x.name}
            </mwc-list-item>
          `
        )}
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
      <div class="select-resource">
        <mwc-button
          raised
          label="Select a resource"
          @click=${this.openResourceOptions}
        ></mwc-button>
        ${this.getResourceOptions()}
      </div>
      <google-chart id="membership-access-trends"> </google-chart>
    `;
  }
}
