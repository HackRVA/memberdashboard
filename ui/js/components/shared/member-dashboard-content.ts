// lit element
import {
  CSSResult,
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// material
import "@material/mwc-top-app-bar-fixed";
import "@material/mwc-tab-bar";
import "@material/mwc-tab";
import "@material/mwc-icon-button";
import "@material/mwc-menu";

// vaadin
import { Router } from "@vaadin/router";

// membership
import { TabIndex } from "../../enums";
import { AuthService } from "../../service";
import { memberDashboardContentStyles } from "./styles";

@customElement("member-dashboard-content")
export class MemberDashboardContent extends LitElement {
  @property({ type: String })
  email: string;

  authService: AuthService = new AuthService();

  static get styles(): CSSResult[] {
    return [memberDashboardContentStyles];
  }
  goToHome(): void {
    Router.go("/home");
  }

  goToUser(): void {
    Router.go("/user");
  }

  goToReports(): void {
    Router.go("/reports");
  }

  goToMembers(): void {
    Router.go("/members");
  }

  goToResources(): void {
    Router.go("/resources");
  }

  getTabIndex(pathName: string): number {
    switch (pathName) {
      case "/home":
        return TabIndex.home;
      case "/user":
        return TabIndex.user;
      case "/reports":
        return TabIndex.reports;
      case "/members":
        return TabIndex.members;
      case "/resources":
        return TabIndex.resources;
      default:
        return -1;
    }
  }

  handleLogout(): void {
    this.authService.logout().subscribe({
      next: (response: null) => {
        localStorage.removeItem("jwt");
        window.location.reload();
      },
    });
  }

  displayLogout(): TemplateResult {
    return html`
      <mwc-button
        class="logout"
        slot="actionItems"
        label="Log out"
        icon="logout"
        @click=${this.handleLogout}
      ></mwc-button>
    `;
  }

  render(): TemplateResult {
    return html`
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
        <div slot="actionItems">${this.email}</div>
        ${this.displayLogout()}
      </mwc-top-app-bar-fixed>
      <mwc-tab-bar activeIndex=${this.getTabIndex(window.location.pathname)}>
        <mwc-tab label="Home" @click=${this.goToHome}></mwc-tab>
        <mwc-tab label="User" @click=${this.goToUser}></mwc-tab>
        <mwc-tab label="Reports" @click=${this.goToReports}></mwc-tab>
        <mwc-tab label="Members" @click=${this.goToMembers}></mwc-tab>
        <mwc-tab label="Resources" @click=${this.goToResources}></mwc-tab>
      </mwc-tab-bar>

      <slot> </slot>
    `;
  }
}
