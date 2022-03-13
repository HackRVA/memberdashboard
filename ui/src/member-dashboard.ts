// lit element
import { customElement } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

// vaadin
import { Router, RouterLocation } from '@vaadin/router';

// memberdashboard
import './router';
import './auth/components/login-page';
import './shared/components/md-content';
import './shared/components/loading-content';
import './material-loader';
import { UserResponse } from './user/types/api/user-response';
import { UserService } from './user/services/user.service';
import { authUser$ } from './auth/auth-user';
import { Inject } from './shared/di';

@customElement('member-dashboard')
export class MemberDashboard extends LitElement {
  email: string;

  @Inject('user')
  private userService: UserService;

  finishedLoading: boolean = false;

  constructor() {
    super();
    // initialize user profile before the app fully loads
    this.getUser();
  }

  onBeforeEnter(location: RouterLocation): void {
    if (location.pathname === '/') {
      this.goToHome();
    }
  }

  goToHome(): void {
    Router.go('/home');
  }

  getUser(): void {
    this.userService
      .getUser()
      .toPromise()
      .then((response: UserResponse) => {
        const { email } = response;
        authUser$.next({ login: true, email: email });
        this.email = email;
        this.finishedLoading = true;
        this.requestUpdate();
      })
      .catch(() => {
        this.finishedLoading = true;
      });
  }

  isUserLogin(): boolean {
    return authUser$.getValue().login;
  }

  loadingAppContent(): TemplateResult {
    return html`
      <loading-content .finishedLoading=${this.finishedLoading}>
        ${this.displayAppContent()}
      </loading-content>
    `;
  }

  displayAppContent(): TemplateResult {
    if (this.isUserLogin()) {
      return html`
        <md-content .email=${this.email}>
          <slot></slot>
        </md-content>
      `;
    } else {
      return html`<login-page></login-page>`;
    }
  }

  render(): TemplateResult {
    return html`${this.loadingAppContent()}`;
  }
}
