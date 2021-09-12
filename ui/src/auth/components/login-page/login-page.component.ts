import { coreStyle } from '../../../shared/styles';

// memberdashboard
import '../login-form';
import { loginPageStyle } from './login-page.style';
import '../register-form';
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

@customElement('login-page')
export class LoginPage extends LitElement {
  isRegister: boolean = false;

  static get styles(): CSSResult[] {
    return [loginPageStyle, coreStyle];
  }

  handleSwitch(): void {
    this.isRegister = !this.isRegister;
    this.requestUpdate();
  }

  displayRegisterLoginForm(): TemplateResult {
    if (this.isRegister) {
      return html`<register-form></register-form>`;
    } else {
      return html`<login-form></login-form>`;
    }
  }

  toggleInfoText(): TemplateResult {
    if (!this.isRegister) {
      return html`
        <span>
          Are you new? Register
          <a href="" @click=${this.handleSwitch}> here </a>
        </span>
      `;
    } else {
      return html`
        <span>
          Already a member? Rock on! Login
          <a href="" @click=${this.handleSwitch}> here </a>
        </span>
      `;
    }
  }

  displayLoginHeaderText(): string {
    return this.isRegister ? 'Register' : 'Login';
  }
  render(): TemplateResult {
    return html`
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
      </mwc-top-app-bar-fixed>
      <h1 class="text-center">${this.displayLoginHeaderText()}</h1>
      <div class="login-container">
        ${this.displayRegisterLoginForm()}
        <div class="toggle-form-text text-center">${this.toggleInfoText()}</div>
      </div>
    `;
  }
}
