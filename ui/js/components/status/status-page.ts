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
import "../shared/card-element";
@customElement("status-page")
export class StatusPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html`
      <card-element>
        <h1>Status</h1>
      </card-element>
    `;
  }
}
