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

@customElement("member-dashboard")
export class MemberDashboard extends LitElement {
  email: string;
  userService: UserService = new UserService();

  onBeforeEnter(location: RouterLocation): void {
    if (location.pathname === "/") {
      this.goToHome();
    }
  }

  goToHome(): void {
    Router.go("/home");
  }

  firstUpdated(): void {
    this.getUser();
  }

  getUser(): void {
    this.userService.getUser().subscribe({
      next: (result: UserProfile) => {
        const { email } = result;
        this.email = email;
        this.requestUpdate();
      },
    });
  }

  isUserLogin(): boolean {
    return !!this.email;
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
