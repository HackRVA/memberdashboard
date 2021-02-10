import { TabIndex } from "./enums";
import { LitElement, html, TemplateResult, customElement } from "lit-element";
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

  displayUserProfile(): TemplateResult {
    if (this.showUserProfile) {
      return html` <user-login-profile .email=${this.email} /> `;
    }
    return html` <login-form /> `;
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

  handleProfileClick(): void {
    const profileBtn: HTMLElement = this.shadowRoot?.querySelector(
      "#profileBtn"
    );
    const menu: HTMLElement & {
      anchor: HTMLElement;
      show: Function;
    } = this.shadowRoot?.querySelector("#menu");

    menu.anchor = profileBtn;
    menu.show();
  }

  render(): TemplateResult {
    return html` <div>
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
        <div slot="actionItems">${this.email}</div>
        <mwc-icon-button
          @click=${this.handleProfileClick}
          id="profileBtn"
          icon="person"
          slot="actionItems"
        ></mwc-icon-button>
        <mwc-menu id="menu" activatable>
          ${this.displayUserProfile()}
        </mwc-menu>
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
    </div>`;
  }
}
