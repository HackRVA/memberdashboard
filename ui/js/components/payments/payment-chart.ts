// lit element
import { LitElement, html, customElement, css, CSSResult } from "lit-element";

// google
import "@google-web-components/google-chart";
import { GoogleChart } from "@google-web-components/google-chart";

// membership
import { PaymentService } from "../../service/payment.service";
import "../shared/card-element";
import { PaymentChartResponse, ChartAttributes } from "./types";

@customElement("payment-chart")
export class NewElement extends LitElement {
  paymentService: PaymentService = new PaymentService();
  paymentCharts: Array<PaymentChartResponse> = [];
  static get styles(): CSSResult {
    return css`
      #payment-chart-container {
        display: flex;
        justify-content: center;
        padding: 36px;
      }
    `;
  }

  firstUpdated() {
    this.handleGetResources();
  }

  addChart(chartData: PaymentChartResponse) {
    const chartContainer = this.shadowRoot?.querySelector("#chart-container");
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

    chartContainer?.appendChild(newChart);
  }

  handleGetResources(): void {
    this.paymentService.getPaymentCharts().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          console.error("some error getting resources");
        } else {
          this.paymentCharts = result as PaymentChartResponse[];

          this.paymentCharts.forEach((x: any) => {
            this.addChart(x);
          });
          this.requestUpdate();
        }
      },
    });
  }
  render() {
    return html`
      <div id="payment-chart-container">
        <chart-container id="chart-container"> </chart-container>
      </div>
    `;
  }
}
