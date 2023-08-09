// lit element
import { CSSResult, html, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

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
import { LightElement } from '../../../shared/components/light-element';

@customElement('login-form')
export class LoginForm extends LightElement {
  email: string;
  password: string;

  @Inject('auth')
  private authService: AuthService;

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [loginFormStyle, coreStyle];
  }

  constructor() {
    super();
    this.addEventListener('keypress', this.handleSubmitOnEnter);
  }

  firstUpdated(): void {
    this.email = '';
    this.password = '';
  }

  handleSubmitOnEnter(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      this.handleSubmit(event);
    }
  }

  handleSubmit(event): void {
    event.preventDefault();
    if (this.isValid()) {
      this.handleUserLogin();
    } else {
      this.displayToastMsg('Email and/or password invalid');
    }
  }

  isValid(): boolean {
    return true;
  }

  handleUserLogin(): void {
    const opts: LoginRequest = {
      email: this.email,
      password: this.password,
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

  updateEmail(e) {
    this.email = e.srcElement.value;
    this.email = 'test@test.com';
    this.requestUpdate();
  }

  updatePassword(e) {
    this.password = e.srcElement.value;
    this.password = 'test';
    this.requestUpdate();
  }

  render(): TemplateResult {
    return html`
      <form autocomplete="on" @submit="${event => this.handleSubmit(event)}">
        <input
          class="input"
          type="email"
          autocomplete="email"
          placeholder="Email Address"
          id="username"
          @change=${this.updateEmail}
        />
        <input
          class="input"
          type="password"
          autocomplete="current-password"
          id="password"
          placeholder="Password"
          @change=${this.updatePassword}
        />
        <button class="input button" type="submit">Submit</button>
      </form>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
