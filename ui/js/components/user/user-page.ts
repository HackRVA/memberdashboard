import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

@customElement("user-page")
export class UserPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html` <h1>User</h1> `;
  }
}
