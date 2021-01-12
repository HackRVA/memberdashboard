import { LitElement, html, TemplateResult, customElement } from "lit-element";
import '@material/mwc-tab-bar';
import '@material/mwc-tab';
import "./components/top-bar";
import { Router } from '@vaadin/router'

@customElement("member-dashboard")
export class MemberDashboard extends LitElement {

  goToHome(): void {
    Router.go('/build')
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

  render(): TemplateResult {
    
    return html`
    <div>
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">
          Member Dashboard
        </div>
        <mwc-icon-button
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
      <div> HI</div>

      <slot> </slot>
    </div>`;
  }
}
