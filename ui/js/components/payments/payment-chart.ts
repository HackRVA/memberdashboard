// lit element
import {
  LitElement,
  html,
  customElement,
  CSSResult,
  TemplateResult,
} from "lit-element";

// google
import "@google-web-components/google-chart";
import { GoogleChart } from "@google-web-components/google-chart";

// membership
import { PaymentService } from "../../service/payment.service";
import "../shared/card-element";
import { PaymentChartResponse, ChartAttributes } from "./types";
import { paymentChartStyles } from "./styles/payment-chart-styles";
@customElement("payment-chart")
export class PaymentChart extends LitElement {
  paymentService: PaymentService = new PaymentService();
  paymentCharts: Array<PaymentChartResponse> = [];
  chartContainer: Element;

  static get styles(): CSSResult[] {
    return [paymentChartStyles];
  }

  firstUpdated(): void {
    this.chartContainer = this.shadowRoot?.querySelector("#chart-container");
    this.createPaymentReportChart();
  }

  addChart(chartData: PaymentChartResponse): void {
    const newChartAttributes: ChartAttributes = {
      id: chartData.id,
      type: chartData.type,
      options: JSON.stringify(chartData.options),
      rows: JSON.stringify(chartData.rows),
      cols: JSON.stringify(chartData.cols),
    };

    const newChart: GoogleChart = document.createElement("google-chart");

    let key: keyof typeof newChartAttributes;
    for (key in newChartAttributes) {
      newChart.setAttribute(key, newChartAttributes[key]);
    }

    this.chartContainer?.appendChild(newChart);
  }

  createPaymentReportChart(): void {
    this.paymentService.getPaymentCharts().subscribe({
      next: (result: any) => {
        this.paymentCharts = result as PaymentChartResponse[];
        this.paymentCharts.forEach((x: PaymentChartResponse) => {
          this.addChart(x);
        });
        this.requestUpdate();
      },
      error: () => {
        console.error("unable to create payment report chart");
      },
    });
  }

  render(): TemplateResult {
    return html`
      <div id="payment-chart-container">
        <chart-container id="chart-container"> </chart-container>
      </div>
    `;
  }
}
