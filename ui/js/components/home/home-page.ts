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
import { UserService } from "../../service/user.service";
import "../shared/card-element";
import "../shared/register-form";

@customElement("home-page")
export class HomePage extends LitElement {
  userService: UserService = new UserService();

  static get styles(): CSSResult {
    return css`
      login-container {
        display: grid;
        justify-content: center;
        padding: 24px;
      }

      .center {
        text-align: center;
      }
    `;
  }

  isUserLogin(): boolean {
    return !!localStorage.getItem("jwt");
  }

  displayHomePage(): TemplateResult {
    if (this.isUserLogin()) {
      return html` <h1>Home</h1> `;
    } else {
      return html`
        <div>
          <h1>Home</h1>
          <login-container>
            <register-form />
          </login-container>
        </div>
      `;
    }
  }

  render(): TemplateResult {
    return html` <card-element> ${this.displayHomePage()} </card-element> `;
  }
}
