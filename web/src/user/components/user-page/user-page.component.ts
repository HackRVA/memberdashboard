// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import './../../../shared/components/md-card';
import { authUser$ } from '../../../auth/auth-user';
import '../user-detail';

@customElement('user-page')
export class UserPage extends LitElement {
  @property({ type: String })
  email: string = '';

  static get styles(): CSSResult[] {
    return [];
  }

  firstUpdated(): void {
    this.email = authUser$.getValue().email;
  }

  displayUserDetail(): TemplateResult | void {
    if (this.email) {
      return html` <user-detail .email=${this.email}> </user-detail> `;
    }
  }

  render(): TemplateResult {
    return html` <md-card> ${this.displayUserDetail()} </md-card> `;
  }
}