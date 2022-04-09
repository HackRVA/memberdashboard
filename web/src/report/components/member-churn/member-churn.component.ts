// lit element
import { customElement } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

// polymer
import '@polymer/paper-tooltip';

// memberdashboard
import {
  primaryDarkGray,
  primaryDarkGreen,
  primaryRed,
} from '../../../shared/styles/colors';
import { ReportService } from '../../services/report.service';
import { ChurnResponse } from '../../types/api/churn-response';
import { Inject } from '../../../shared/di/inject';

@customElement('member-churn')
export class MemberChurn extends LitElement {
  @Inject('report')
  private reportService: ReportService;

  churn: number;

  firstUpdated(): void {
    this.getMemberChurn();
  }

  getMemberChurn() {
    this.reportService.getMemberChurn().subscribe({
      next: (response: ChurnResponse) => {
        this.churn = response.churn;
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to get member churn');
      },
    });
  }

  churnTemplate(): TemplateResult {
    if (this.churn < 0) {
      return html` <span style="color:${primaryRed}">${this.churn}</span> `;
    }
    if (this.churn === 0) {
      return html`
        <span style="color:${primaryDarkGray}">${this.churn}</span>
      `;
    }
    return html` <span style="color:${primaryDarkGreen}">${this.churn}</span> `;
  }

  render(): TemplateResult {
    return html`
      <div>
        <paper-tooltip for="member-churn" animation-delay="0">
          Churn is how many members we have lost or gained since last month.
        </paper-tooltip>
        <h3 id="member-churn">Member Churn: ${this.churnTemplate()}</h3>
      </div>
    `;
  }
}
