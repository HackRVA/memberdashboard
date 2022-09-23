// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement } from 'lit/decorators.js';

// vaadin
import { Router } from '@vaadin/router';

// memberdashboard
import { notFoundStyle } from './not-found.style';
import { coreStyle } from '../../styles';

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
        <div class="not-found-container text-center">
          <h2> <b>404</b> </h2>
          <p class="sad-freshmon">${'(っ- ‸ – ς)'}</p>
          <p class="page-not-found">Page not found</p>
          <mwc-button label="Go back to the home page" dense unelevated @click=${
            this.goBackToHomePage
          }> </mvc-button>
        </div>
      `;
  }
}
