// lit element
import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import "./payment-chart";
import "../shared/card-element";

@customElement("payments-page")
export class PaymentsPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html`
      <card-element>
        <payment-chart> </payment-chart>
      </card-element>
    `;
  }
}
