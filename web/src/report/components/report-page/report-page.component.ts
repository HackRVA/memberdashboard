// lit element
import { customElement } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../report-chart';
import '../member-churn';
import '../../../shared/components/md-card';
import { withCard } from '../../../shared/functions';

@customElement('report-page')
export class ReportPage extends LitElement {
  render(): TemplateResult {
    return withCard(html`
      <member-churn> </member-churn>
      <report-chart> </report-chart>
    `);
  }
}
