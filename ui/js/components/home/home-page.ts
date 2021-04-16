// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import { homePageStyles } from "./styles/home-page-styles";
import "../shared/card-element";

@customElement("home-page")
export class HomePage extends LitElement {
  static get styles(): CSSResult[] {
    return [homePageStyles];
  }

  firstUpdated(): void {}

  displayHomePage(): TemplateResult {
    return html``;
  }

  render(): TemplateResult {
    return html` <card-element> ${this.displayHomePage()} </card-element> `;
  }
}
