import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import "./payment-chart";

@customElement("payments-page")
export class PaymentsPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html` <payment-chart> </payment-chart> `;
  }
}
