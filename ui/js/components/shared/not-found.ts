// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// vaadin
import { Router } from "@vaadin/router";

// membership
import "../shared/card-element";
import { notFoundStyles } from "./styles";

@customElement("not-found")
export class NotFound extends LitElement {
  static get styles(): CSSResult[] {
    return [notFoundStyles];
  }

  goBackToHomePage(): void {
    Router.go("/home");
  }

  render(): TemplateResult {
    return html` 
      <div class="not-found-container text-center">
        <div>
          <strong>404</strong>
        </div>
        <div class="sad-freshmon">${"(っ- ‸ – ς)"}</div>
        <div class="page-not-found">Page not found</div>
        <mwc-button label="Go back to the home page" dense unelevated @click=${
          this.goBackToHomePage
        }> </mvc-button>
      </div>
    `;
  }
}
