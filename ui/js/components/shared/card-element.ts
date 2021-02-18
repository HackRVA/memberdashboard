// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import { cardElementStyles } from "./styles";

@customElement("card-element")
export class CardElement extends LitElement {
  static get styles(): CSSResult[] {
    return [cardElementStyles];
  }
  render(): TemplateResult {
    return html`
      <card-container>
        <div class="card">
          <div class="container">
            <slot></slot>
          </div>
        </div>
      </card-container>
    `;
  }
}
