import { TabIndex } from "./enums";
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
  css,
} from "lit-element";
import "@material/mwc-tab-bar";
import "@material/mwc-tab";
import "@material/mwc-top-app-bar-fixed";
import "@material/mwc-icon-button";
import "@material/mwc-menu";
import "./components/shared/login-form";
import "./components/user/user-login-profile";
import "./router";
import { Router, RouterLocation } from "@vaadin/router";
import { UserService } from "./service/user.service";

@customElement("member-dashboard")
export class MemberDashboard extends LitElement {
  showUserProfile: boolean = false;
  email: string = "";
  userService: UserService = new UserService();

  static get styles(): CSSResult {
    return css`
      .logout {
        margin-left: 8px;
        --mdc-theme-primary: white;
      }
    `;
  }

  onBeforeEnter(location: RouterLocation): void {
    if (location.pathname === "/") {
      Router.go("/home");
    }
  }

  goToHome(): void {
    Router.go("/home");
  }

  goToUser(): void {
    Router.go("/user");
  }

  goToPayments(): void {
    Router.go("/payments");
  }

  goToMembers(): void {
    Router.go("/members");
  }

  goToResources(): void {
    Router.go("/resources");
  }

  goToStatus(): void {
    Router.go("/status");
  }

  updated(): void {
    if (this.showUserProfile) return;

    this.userService.getUser().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        const { email } = result as UserService.UserProfile;
        this.email = email;
        this.showUserProfile = true;
        this.requestUpdate();
      },
    });
  }

  handleLogout(): void {
    this.userService.logout().subscribe({
      next: (response: null) => {
        localStorage.removeItem("jwt");
        window.location.reload();
      },
    });
  }

  isUserLogin(): boolean {
    return !!localStorage.getItem("jwt");
  }

  displayLogout(): TemplateResult {
    if (this.isUserLogin()) {
      return html`
        <mwc-button
          class="logout"
          slot="actionItems"
          label="Log out"
          @click=${this.handleLogout}
        ></mwc-button>
      `;
    }
    return html``;
  }

  getTabIndex(pathName: string): number {
    switch (pathName) {
      case "/home":
        return TabIndex.home;
      case "/user":
        return TabIndex.user;
      case "/payments":
        return TabIndex.payments;
      case "/members":
        return TabIndex.members;
      case "/resources":
        return TabIndex.resources;
      case "/status":
        return TabIndex.status;
      default:
        return -1;
    }
  }

  render(): TemplateResult {
    return html`
      <div>
        <mwc-top-app-bar-fixed centerTitle>
          <div slot="title">Member Dashboard</div>
          <div slot="actionItems">${this.email}</div>
          ${this.displayLogout()}
        </mwc-top-app-bar-fixed>
        <mwc-tab-bar activeIndex=${this.getTabIndex(window.location.pathname)}>
          <mwc-tab label="Home" @click=${this.goToHome}></mwc-tab>
          <mwc-tab label="User" @click=${this.goToUser}></mwc-tab>
          <mwc-tab label="Payments" @click=${this.goToPayments}></mwc-tab>
          <mwc-tab label="Members" @click=${this.goToMembers}></mwc-tab>
          <mwc-tab label="Resources" @click=${this.goToResources}></mwc-tab>
          <mwc-tab label="Status" @click=${this.goToStatus}></mwc-tab>
        </mwc-tab-bar>

        <slot> </slot>
      </div>
    `;
  }
}
