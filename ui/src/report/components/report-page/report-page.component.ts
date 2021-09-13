// lit element
import { customElement } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../payment-chart';
import '../../../shared/components/md-card';

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
