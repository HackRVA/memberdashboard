// lit element
import { LitElement, html, customElement, TemplateResult } from "lit-element";

// membership
import "./payment-chart";
import "../shared/card-element";

@customElement("payments-page")
export class PaymentsPage extends LitElement {
  render(): TemplateResult {
    return html`
      <card-element>
        <payment-chart> </payment-chart>
      </card-element>
    `;
  }
}
