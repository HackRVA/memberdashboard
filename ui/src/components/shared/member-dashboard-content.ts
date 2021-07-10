// lit element
import {
  CSSResult,
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// vaadin
import { Router } from "@vaadin/router";

// memberdashboard
import { TabIndex } from "../../enums";
import { AuthService, VersionService } from "../../service";
import { memberDashboardContentStyles, coreStyles } from "./styles";
import { isAdmin } from "./../../function";
import { VersionResponse } from "./types";

@customElement("member-dashboard-content")
export class MemberDashboardContent extends LitElement {
  @property({ type: String })
  email: string;

  version: VersionResponse;

  authService: AuthService = new AuthService();
  versionService: VersionService = new VersionService();

  static get styles(): CSSResult[] {
    return [memberDashboardContentStyles, coreStyles];
  }

  firstUpdated(): void {
    this.getVersion();
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
        window.location.href = "/home";
      },
    });
  }

  getVersion(): void {
    this.versionService.getVersion().subscribe({
      next: (response: VersionResponse) => {
        this.version = response;
        this.requestUpdate();
      },
    });
  }

  generateVersionNumber(version: VersionResponse): TemplateResult {
    return html`
      <span>
        ${version?.major}.${version?.minor}.${version?.hotfix}.<b
          >${version?.build}</b
        >
      </span>
    `;
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

  renderAdminTabs(): TemplateResult | void {
    if (isAdmin()) {
      return html`
        <mwc-tab
          label="Reports"
          icon="show_chart"
          @click=${this.goToReports}
        ></mwc-tab>
        <mwc-tab
          label="Members"
          icon="people"
          @click=${this.goToMembers}
        ></mwc-tab>
        <mwc-tab
          label="Resources"
          icon="devices"
          @click=${this.goToResources}
        ></mwc-tab>
      `;
    }

    return html``;
  }

  render(): TemplateResult {
    return html`
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
        <div slot="actionItems">${this.email}</div>
        ${this.displayLogout()}
      </mwc-top-app-bar-fixed>
      <mwc-tab-bar activeIndex=${this.getTabIndex(window.location.pathname)}>
        <mwc-tab label="Home" icon="home" @click=${this.goToHome}></mwc-tab>
        <mwc-tab
          label="User"
          icon="account_circle"
          @click=${this.goToUser}
        ></mwc-tab>
        ${this.renderAdminTabs()}
      </mwc-tab-bar>

      <slot> </slot>

      <div class="version margin-r-24">
        <p>Version ${this.generateVersionNumber(this.version)}</p>
      </div>
    `;
  }
}
