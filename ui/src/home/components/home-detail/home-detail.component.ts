// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement } from 'lit/decorators.js';

// memberdashboard
import { coreStyle } from './../../../shared/styles/core.style';

@customElement('home-detail')
export class HomeDetail extends LitElement {
  static get styles(): CSSResult[] {
    return [coreStyle];
  }

  render(): TemplateResult {
    return html`<div>HELLO WORLD</div>`;
  }
}
