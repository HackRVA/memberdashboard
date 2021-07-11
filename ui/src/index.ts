// lit element
import { LitElement, html, TemplateResult, customElement } from "lit-element";

// vaadin
import { Router, RouterLocation } from "@vaadin/router";

// material
import "./material-loader";

// memberdashboard
import "./router";
import { UserService } from "./service/user.service";
import { UserProfile } from "./components/user/types";
import "./components/shared/login-page";
import "./components/shared/member-dashboard-content";
import { authUser } from "./auth-user";

@customElement("member-dashboard")
export class MemberDashboard extends LitElement {
  email: string;
  userService: UserService = new UserService();

  constructor() {
    super();
    // initialize user profile before the app fully loads
    this.getUser();
  }

  onBeforeEnter(location: RouterLocation): void {
    if (location.pathname === "/") {
      this.goToHome();
    }
  }

  goToHome(): void {
    Router.go("/home");
  }

  getUser(): void {
    this.userService.getUser().subscribe({
      next: (result: UserProfile) => {
        const { email } = result;
        authUser.next({ login: true, email: email });
        this.email = email;
        this.requestUpdate();
      },
    });
  }

  isUserLogin(): boolean {
    return authUser.getValue().login;
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
