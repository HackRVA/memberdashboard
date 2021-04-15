// lit element
import {
  css,
  CSSResult,
  customElement,
  html,
  LitElement,
  TemplateResult,
} from "lit-element";

// membership
import "./card-element";
import "./login-form";
import "./register-form";

@customElement("login-page")
export class LoginPage extends LitElement {
  isRegister: boolean = false;

  static get styles(): CSSResult {
    return css`
      .login-container {
        display: grid;
        justify-content: center;
      }

      .text-center {
        text-align: center;
      }

      h1 {
        margin-top: 0px;
      }
    `;
  }

  handleSwitch(): void {
    this.isRegister = !this.isRegister;
    this.requestUpdate();
  }

  displayRegisterLoginForm(): TemplateResult {
    if (this.isRegister) {
      return html`<register-form @switch=${this.handleSwitch} />`;
    } else {
      return html`<login-form @switch=${this.handleSwitch} />`;
    }
  }

  displayLoginHeaderText(): string {
    return this.isRegister ? "Register" : "Login";
  }
  render(): TemplateResult {
    return html`
      <card-element class="text-center">
        <h1>${this.displayLoginHeaderText()}</h1>
        <div class="login-container">${this.displayRegisterLoginForm()}</div>
      </card-element>
    `;
  }
}
