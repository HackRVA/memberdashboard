import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

@customElement("status-page")
export class StatusPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html` <h1>Status</h1> `;
  }
}
