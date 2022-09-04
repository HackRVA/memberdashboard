// lit element
import { html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// memberdashboard
import { Inject } from '../../../../shared/di';
import { Member } from '../../../types/api/member-response';
import { MemberService } from '../../../services/member.service';

@customElement('refresh-member-status-form')
export class RemoveResourceForm extends LitElement {
  @property({ type: String })
  subscriptionID: string = '';

  @property({ type: Function })
  closeHandler: () => void;

  @Inject('member')
  private memberService: MemberService;

  member: Member;
  errorMessage: string;

  private checkMemeberLevel(): void {
    this.memberService
      .checkMemberStatus(this.subscriptionID)
      .toPromise()
      .then(result => (this.member = result))
      .catch(err => (this.errorMessage = err.message))
      .finally(() => this.requestUpdate());
  }

  async firstUpdated(): Promise<void> {
    this.checkMemeberLevel();
  }

  private getOutputMsg() {
    if (this.errorMessage) {
      return this.errorMessage;
    }

    return this.member.memberLevel;
  }
  render(): TemplateResult {
    return html`
      <p>${this.subscriptionID}</p>
      <p>${!this.member ? 'fetching...' : this.getOutputMsg()}</p>
      <mwc-button slot="secondaryAction" @click=${this.closeHandler}>
        Close
      </mwc-button>
    `;
  }
}
