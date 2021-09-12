import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('member-dashboard-content')
export class MemberDashboardContent extends LitElement {
  render() {
    return html`
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
        <div slot="actionItems"></div>
      </mwc-top-app-bar-fixed>
      <mwc-tab-bar>
        <mwc-tab label="Home" icon="home"></mwc-tab>
        <mwc-tab label="User" icon="account_circle"></mwc-tab>
      </mwc-tab-bar>
    `;
  }
}
