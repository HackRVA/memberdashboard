import { LitElement, html, TemplateResult, customElement } from "lit-element";
import '@material/mwc-tab-bar';
import '@material/mwc-tab';
import "./components/top-bar";
import "./components/login-form";
import "./components/user-login-profile";
import { Router, RouterLocation } from '@vaadin/router'
import { UserService } from "./service/User";

@customElement("member-dashboard")
export class MemberDashboard extends LitElement {
  showUserProfile: boolean = false;
  username: string = "";
  email: string = "";
  userService: UserService = new UserService();

  onBeforeEnter(location: RouterLocation): void {
    if (location.pathname === '/') {
      Router.go('/home');
    }
  }

  goToHome(): void {
    Router.go('/home')
  }

  goToUsers(): void {
    Router.go('/users')
  }

  goToMembers(): void {
    Router.go('/members')
  }

  goToResources(): void {
    Router.go('/resources')
  }

  goToStatus(): void {
    Router.go('/status')
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
        const { username, email } = result as UserService.UserProfile;
        this.username = username;
        this.email = email;
        this.showUserProfile = true;
        this.requestUpdate();
      },
    });
  }

  displayUserProfile(): TemplateResult | undefined {
    if (this.showUserProfile) {
      return html`
        <user-login-profile />
      `
    }
    return html `
      <login-form />
    `
  }

  handleProfileClick(): void {
    const profileBtn:
      | HTMLElement
      | null
      | undefined = this.shadowRoot?.querySelector("#profileBtn");
    const menu:
      | (HTMLElement & { anchor: HTMLElement; show: Function })
      | null
      | undefined = this.shadowRoot?.querySelector("#menu");

    if (!profileBtn) return console.error("profile btn doesn't exist");
    if (!menu) return console.error("menu element doesn't exist");

    menu.anchor = profileBtn;
    menu.show();
  }

  render(): TemplateResult {
    
    return html`
    <div>
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">
          Member Dashboard
        </div>
        <div slot="actionItems"></div>
        <mwc-icon-button
          @click=${this.handleProfileClick}
          id="profileBtn"
          icon="person"
          slot="actionItems"
        ></mwc-icon-button>
        <mwc-menu id="menu" activatable> ${this.displayUserProfile()} </mwc-menu>
      </mwc-top-app-bar-fixed>
      <mwc-tab-bar>
        <mwc-tab label="Home" @click=${this.goToHome}></mwc-tab>
        <mwc-tab label="Users" @click=${this.goToUsers}></mwc-tab>
        <mwc-tab label="Members" @click=${this.goToMembers}></mwc-tab>
        <mwc-tab label="Resources" @click=${this.goToResources}></mwc-tab>
        <mwc-tab label="Status" @click=${this.goToStatus}></mwc-tab>
      </mwc-tab-bar>

      <slot> </slot>
    </div>`;
  }
}
