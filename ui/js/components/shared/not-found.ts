import { css } from "lit-element";
import "@material/mwc-button";
// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import "../shared/card-element";
import { Router } from "@vaadin/router";

// membership

@customElement("not-found")
export class NotFound extends LitElement {
  static get styles(): CSSResult {
    return css`
      .not-found-container {
        font-size: 36px;
      }
      .text-center {
        text-align: center;
      }
      .sad-freshmon {
        margin-top: 32px;
        margin-bottom: 32px;
        font-size: 70px;
      }
      .page-not-found {
        margin-bottom: 32px;
      }
    `;
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
