import { AuthResponse } from './../../types/api/auth-response';

// material
import { TextField } from '@material/mwc-textfield/mwc-textfield';

// memberdashboard
import { coreStyle } from '../../../shared/styles';
import '../../../shared/components/toast-msg';
import { loginFormStyle } from './login-form.style';
import { LoginRequest } from '../../types/api/login-request';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { AuthService } from '../../services/auth.service';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { showComponent } from '../../../shared/functions';

@customElement('login-form')
export class LoginForm extends LitElement {
  // form template
  emailFieldTemplate: TextField;
  passwordFieldTemplate: TextField;

  authService: AuthService = new AuthService();

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [loginFormStyle, coreStyle];
  }

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot?.querySelector('#email');
    this.passwordFieldTemplate = this.shadowRoot?.querySelector('#password');
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.handleUserLogin();
    } else {
      this.displayToastMsg('Email and/or password invalid');
    }
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.passwordFieldTemplate.validity.valid
    );
  }

  handleUserLogin(): void {
    const opts: LoginRequest = {
      email: this.emailFieldTemplate?.value,
      password: this.passwordFieldTemplate?.value,
    };

    this.authService.login(opts).subscribe({
      next: (result: AuthResponse) => {
        const { token } = result;
        localStorage.setItem('jwt', token);
        window.location.reload();
      },
      error: () => {
        this.displayToastMsg('Email and/or password invalid');
      },
    });
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  render(): TemplateResult {
    return html`
      <mwc-formfield>
        <mwc-textfield
          size="30"
          id="email"
          required
          type="email"
          label="Email"
        ></mwc-textfield>
      </mwc-formfield>
      <mwc-formfield>
        <mwc-textfield
          size="30"
          id="password"
          required
          type="password"
          label="Password"
        ></mwc-textfield>
      </mwc-formfield>
      <a href=""> Forgot Password? </a>
      <mwc-button
        unelevated
        label="login"
        @click=${this.handleSubmit}
        class="login"
      ></mwc-button>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}