// lit element

// memberdashboard
import '../payment-chart';
import '../../../shared/components/md-card';
import { customElement } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

@customElement('report-page')
export class ReportPage extends LitElement {
  render(): TemplateResult {
    return html`
      <md-card>
        <payment-chart> </payment-chart>
      </md-card>
    `;
  }
}
