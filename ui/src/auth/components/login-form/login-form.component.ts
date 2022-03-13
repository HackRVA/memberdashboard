// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement } from 'lit/decorators.js';

// material
import { TextField } from '@material/mwc-textfield/mwc-textfield';

// memberdashboard
import '../../../shared/components/toast-msg';
import { AuthResponse } from './../../types/api/auth-response';
import { coreStyle } from '../../../shared/styles';
import { loginFormStyle } from './login-form.style';
import { LoginRequest } from '../../types/api/login-request';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { AuthService } from '../../services/auth.service';
import { showComponent } from '../../../shared/functions';
import { Inject } from '../../../shared/di/inject';

@customElement('login-form')
export class LoginForm extends LitElement {
  // form template
  emailFieldTemplate: TextField;
  passwordFieldTemplate: TextField;

  @Inject('auth')
  private authService: AuthService;

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [loginFormStyle, coreStyle];
  }

  constructor() {
    super();
    this.addEventListener('keypress', this.handleSubmitByEnter);
  }

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot?.querySelector('#email');
    this.passwordFieldTemplate = this.shadowRoot?.querySelector('#password');
  }

  handleSubmitByEnter(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      this.handleSubmit();
    }
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
