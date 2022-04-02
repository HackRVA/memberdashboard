// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// memberdashboard
import { coreStyle } from '../../styles';

@customElement('loading-content')
export class LoadingContent extends LitElement {
  @property({ type: Boolean })
  finishedLoading: boolean;

  static get styles(): CSSResult[] {
    return [coreStyle];
  }

  displayPage(): TemplateResult {
    if (this.finishedLoading) {
      return html` <slot> </slot> `;
    }

    return html`
      <div class="center-text">
        <mwc-circular-progress indeterminate></mwc-circular-progress>
      </div>
    `;
  }

  render(): TemplateResult {
    return html` ${this.displayPage()} `;
  }
}
