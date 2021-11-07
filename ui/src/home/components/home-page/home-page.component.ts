import { isAdmin } from './../../../shared/functions';
// lit element
import { customElement } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../../../shared/components/md-card';
import { homePageStyle } from './home-page.style';
import '../../../shared/components/happy-minion';
import '../home-detail';

@customElement('home-page')
export class HomePage extends LitElement {
  static get styles(): CSSResult[] {
    return [homePageStyle];
  }

  displayContentBasedOnPermission(): TemplateResult {
    if (isAdmin()) {
      return html`<home-detail> </home-detail>`;
    }

    return html`<happy-minion></happy-minion>`;
  }

  render(): TemplateResult {
    return html`
      <md-card> ${this.displayContentBasedOnPermission()} </md-card>
    `;
  }
}
