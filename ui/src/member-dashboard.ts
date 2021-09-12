import { customElement, property } from 'lit/decorators.js';
import { LitElement, html, css } from 'lit';

// material
import './material-loader';
import './shared/components/member-dashboard-content/member-dashboard-content';

import './router';

@customElement('member-dashboard')
export class MemberDashboard extends LitElement {
  render() {
    return html` <member-dashboard-content> </member-dashboard-content> `;
  }
}
