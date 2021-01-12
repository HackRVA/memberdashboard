import { LitElement, html, TemplateResult, customElement } from "lit-element";
import '@material/mwc-tab-bar';
import '@material/mwc-tab';
import "./components/top-bar";
import { Router, RouterLocation } from '@vaadin/router'

@customElement("member-dashboard")
export class MemberDashboard extends LitElement {

  onBeforeEnter(location: RouterLocation): void {
    if (location.pathname === '/build/') {
      Router.go('/build/home');
    }
  }

  goToHome(): void {
    Router.go('/build/home')
  }

  goToUsers(): void {
    Router.go('/build/users')
  }

  goToMembers(): void {
    Router.go('/build/members')
  }

  goToResources(): void {
    Router.go('/build/resources')
  }

  goToStatus(): void {
    Router.go('/build/status')
  }

  handleProfileClick(): void {
    return;
    // const profileBtn:
    //   | HTMLElement
    //   | null
    //   | undefined = this.shadowRoot?.querySelector("#profileBtn");
    // const menu:
    //   | (HTMLElement & { anchor: HTMLElement; show: Function })
    //   | null
    //   | undefined = this.shadowRoot?.querySelector("#menu");

    // if (!profileBtn) return console.error("profile btn doesn't exist");
    // if (!menu) return console.error("menu element doesn't exist");

    // menu.anchor = profileBtn;
    // menu.show();
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
