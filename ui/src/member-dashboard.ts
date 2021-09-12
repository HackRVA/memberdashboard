import { UserResponse } from './user/types/api/user-response';
// lit element

// vaadin
import { Router, RouterLocation } from '@vaadin/router';

// material
import './material-loader';

// memberdashboard
import './router';
import './auth/components/login-page';
import './components/shared/member-dashboard-content';
import { UserService } from './user/services/user.service';
import { authUser$ } from './auth/auth-user';
import { customElement } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

@customElement('member-dashboard')
export class MemberDashboard extends LitElement {
  email: string;
  userService: UserService = new UserService();

  constructor() {
    super();
    // initialize user profile before the app fully loads
    // this.getUser();
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
    this.userService.getUser().subscribe({
      next: (result: UserResponse) => {
        const { email } = result;
        authUser$.next({ login: true, email: email });
        this.email = email;
        this.requestUpdate();
      },
    });
  }

  isUserLogin(): boolean {
    return true;
  }

  displayAppContent(): TemplateResult {
    if (this.isUserLogin()) {
      return html`
      <member-dashboard-content .email=${this.email}>
        <slot></slot>
      </member-dashboard-content`;
    } else {
      return html`<login-page></login-page>`;
    }
  }

  render(): TemplateResult {
    return html`${this.displayAppContent()}`;
  }
}
