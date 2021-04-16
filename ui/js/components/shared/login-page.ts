// lit element
import {
  CSSResult,
  customElement,
  html,
  LitElement,
  TemplateResult,
} from "lit-element";

// material
import "@material/mwc-top-app-bar-fixed";

// membership
import "./card-element";
import "./login-form";
import "./register-form";
import { loginPageStyles } from "./styles/";

@customElement("login-page")
export class LoginPage extends LitElement {
  isRegister: boolean = false;

  static get styles(): CSSResult[] {
    return [loginPageStyles];
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
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
      </mwc-top-app-bar-fixed>
      <card-element class="text-center">
        <h1>${this.displayLoginHeaderText()}</h1>
        <div class="login-container">${this.displayRegisterLoginForm()}</div>
      </card-element>
    `;
  }
}
