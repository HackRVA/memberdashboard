import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import "./card-element";
import "../user/user-profile";

@customElement("body-element")
export class BodyElement extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html`
      <main>
        <user-profile></user-profile>
      </main>
    `;
  }
}
