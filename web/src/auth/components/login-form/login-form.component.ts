// lit element
import { html, TemplateResult } from 'lit';
import { customElement } from 'lit/decorators.js';

// memberdashboard
import '../../../shared/components/toast-msg';
import { AuthResponse } from './../../types/api/auth-response';
import { loginFormStyle } from './login-form.style';
import { LoginRequest } from '../../types/api/login-request';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { AuthService } from '../../services/auth.service';
import { showComponent } from '../../../shared/functions';
import { Inject } from '../../../shared/di/inject';
import { LightElement } from '../../../shared/components/light-element';
import { authUser$ } from '../../auth-user';
import { UserService } from '../../../user/services/user.service';
import { UserResponse } from '../../../user/types/api/user-response';

@customElement('login-form')
export class LoginForm extends LightElement {
  email: string;
  password: string;

  @Inject('auth')
  private authService: AuthService;

  @Inject('user')
  private userService: UserService;

  toastMsg: ToastMessage;

  constructor() {
    super();
    this.addEventListener('keypress', this.handleSubmitOnEnter);
  }

  firstUpdated(): void {
    this.email = '';
    this.password = '';
    this.getUser();
    this.requestUpdate();
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

  getUser(): void {
    this.userService
      .getUser()
      .toPromise()
      .then((response: UserResponse) => {
        const { email } = response;
        authUser$.next({ login: true, email: email });
        this.email = email;
        this.requestUpdate();
      })
      .catch(() => {
        this.requestUpdate();
      });
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
    this.requestUpdate();
  }

  updatePassword(e) {
    this.password = e.srcElement.value;
    this.requestUpdate();
  }

  render(): TemplateResult {
    if (authUser$.getValue().login) return html``;

    return html`
      <form autocomplete="on" @submit="${event => this.handleSubmit(event)}">
        <div class="form__group">
          <input
            class="input form__field"
            type="email"
            autocomplete="email"
            placeholder="Email Address"
            id="username"
            @change=${this.updateEmail}
          />
          <label for="Email Address" class="form__label">email</label>
        </div>
        <div class="form__group">
          <input
            class="input form__field"
            type="password"
            autocomplete="current-password"
            id="password"
            placeholder="Password"
            @change=${this.updatePassword}
          />
          <label for="Password" class="form__label">password</label>
        </div>
        <div class="form__group">
          <mwc-button
            label="login"
            unelevated
            @click=${this.handleSubmit}
            type="submit"
          ></mwc-button>
        </div>
      </form>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>

      <style>
        ${loginFormStyle.cssText}
      </style>
    `;
  }
}
