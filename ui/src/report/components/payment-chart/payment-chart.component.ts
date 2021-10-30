// lit element
import { customElement } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// google
import '@google-web-components/google-chart';
import { GoogleChart } from '@google-web-components/google-chart';

// memberdashboard
import {
  primaryBlue,
  primaryDarkGray,
  primaryDarkGreen,
  primaryRed,
} from '../../../shared/styles/colors';
import { showComponent } from '../../../shared/functions';
import { paymentChartStyle } from './payment-chart.style';
import { PaymentService } from '../../services/payment.service';
import {
  ChartOptions,
  PaymentChartResponse,
} from '../../types/api/payment-chart-response';
import { Inject } from '../../../shared/di/inject';

@customElement('payment-chart')
export class PaymentChart extends LitElement {
  @Inject('payment')
  private paymentService: PaymentService;

  paymentCharts: Array<PaymentChartResponse> = [];
  membershipTrendsData: PaymentChartResponse;
  membershipDistributionData: Array<PaymentChartResponse> = [];

  membershipTrendsChart: GoogleChart;
  membershipDistributionChart: GoogleChart;

  static get styles(): CSSResult[] {
    return [paymentChartStyle];
  }

  firstUpdated(): void {
    this.membershipTrendsChart =
      this.shadowRoot?.querySelector('#membership-trends');
    this.membershipDistributionChart = this.shadowRoot?.querySelector(
      '#membership-distribution'
    );
    this.createPaymentReportChart();
  }

  creatMembershipTrendChart(chartData: PaymentChartResponse): void {
    this.membershipTrendsChart.type = chartData.type;
    this.membershipTrendsChart.options = chartData.options;
    this.membershipTrendsChart.rows = chartData.rows;
    this.membershipTrendsChart.cols = chartData.cols;
  }

  createMembershipDistributionChart(chartData: PaymentChartResponse): void {
    const options: ChartOptions = Object.assign({}, chartData.options);
    options.title = options.title + ' - Membership Distribution';

    this.membershipDistributionChart.type = chartData.type;
    this.membershipDistributionChart.options = options;
    this.membershipDistributionChart.rows = chartData.rows;
    this.membershipDistributionChart.cols = chartData.cols;
  }

  openMembershipMonthsOptions(): void {
    this.requestUpdate();
    showComponent('#membership-months', this.shadowRoot);
  }

  createPaymentReportChart(): void {
    this.paymentService.getPaymentCharts().subscribe({
      next: (result: PaymentChartResponse[]) => {
        this.paymentCharts = result;
        this.paymentCharts = this.updatePaymentCharts(this.paymentCharts);

        this.membershipTrendsData = this.paymentCharts[0];
        this.membershipDistributionData = this.paymentCharts.splice(1);

        this.creatMembershipTrendChart(this.membershipTrendsData);
        this.createMembershipDistributionChart(
          this.membershipDistributionData[
            this.membershipDistributionData.length - 1
          ]
        );
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to create payment report chart');
      },
    });
  }

  updatePaymentCharts(
    paymentCharts: PaymentChartResponse[]
  ): PaymentChartResponse[] {
    paymentCharts.forEach((x: PaymentChartResponse) => {
      x.options.colors = [
        primaryBlue.cssText,
        primaryDarkGreen.cssText,
        primaryRed.cssText,
        primaryDarkGray.cssText,
      ];
    });

    return paymentCharts;
  }

  updateMembershipDistributionChart(chartData: PaymentChartResponse): void {
    this.createMembershipDistributionChart(chartData);
  }

  getPaymentMonthOptions(): TemplateResult {
    return html`
      <mwc-menu x="10" y="40" id="membership-months">
        ${this.membershipDistributionData.map((x: PaymentChartResponse) => {
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
      <div class="select-month">
        <mwc-button
          raised
          label="Select a month"
          @click=${this.openMembershipMonthsOptions}
        ></mwc-button>
        ${this.getPaymentMonthOptions()}
      </div>
      <div id="payment-chart-container">
        <google-chart id="membership-distribution"> </google-chart>
        <google-chart id="membership-trends"> </google-chart>
      </div>
    `;
  }
}
