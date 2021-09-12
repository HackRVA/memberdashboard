// vaadin
import { Router } from '@vaadin/router';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement } from 'lit/decorators.js';
import { coreStyle } from '../../styles';

// memberdashboard
import { notFoundStyle } from './not-found.style';

@customElement('not-found')
export class NotFound extends LitElement {
  static get styles(): CSSResult[] {
    return [notFoundStyle, coreStyle];
  }

  goBackToHomePage(): void {
    Router.go('/home');
  }

  render(): TemplateResult {
    return html` 
        <div class="not-found-container center-text">
          <div>
            <strong>404</strong>
          </div>
          <div class="sad-freshmon">${'(っ- ‸ – ς)'}</div>
          <div class="page-not-found">Page not found</div>
          <mwc-button label="Go back to the home page" dense unelevated @click=${
            this.goBackToHomePage
          }> </mvc-button>
        </div>
      `;
  }
}
