import { LitElement, html, customElement, css, CSSResult } from "lit-element";
import "@google-web-components/google-chart";
import { PaymentService } from "../../service/payment.service";
import "../shared/card-element";

@customElement("payment-chart")
export class NewElement extends LitElement {
  paymentService: PaymentService = new PaymentService();
  paymentCharts: Array<PaymentService.PaymentChartResponse> | null = null;
  static get styles(): CSSResult {
    return css`
      #payment-chart-container {
        display: flex;
        justify-content: center;
      }
    `;
  }

  firstUpdated() {
    this.handleGetResources();
  }

  addChart(chartData: PaymentService.PaymentChartResponse) {
    const chartContainer = this.shadowRoot?.querySelector("#chart-container");

    const newChartAttributes: any = {
      id: "new-chart",
      type: chartData.type,
      options: JSON.stringify(chartData.options),
      rows: JSON.stringify(chartData.rows),
      cols: JSON.stringify(chartData.cols),
    };

    var newChart = document.createElement("google-chart");

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
          // this.onLoginComplete("Some error logging in");
          console.error("some error getting resources");
        } else {
          this.paymentCharts = result as PaymentService.PaymentChartResponse[];

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
      <card-element>
        <div id="payment-chart-container">
          <chart-container id="chart-container"> </chart-container>
        </div>
      </card-element>
    `;
  }
}
