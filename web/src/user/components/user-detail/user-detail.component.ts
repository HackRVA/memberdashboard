// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// polymer
import '@polymer/paper-card';

// memberdashboard
import { coreStyle } from '../../../shared/styles';
import { userDetailStyle } from './user-detail.style';
import {
  MemberResource,
  Member,
} from '../../../member/types/api/member-response';
import { MemberService } from '../../../member/services/member.service';
import { displayMemberStatus } from '../../../member/functions';
import { showComponent } from '../../../shared/functions';
import './../../../shared/components/rfid-modal';
import { Inject } from '../../../shared/di';

@customElement('user-detail')
export class UserDetail extends LitElement {
  @property({ type: String })
  email: string = '';
  memberUser: Member;

  @Inject('member')
  private memberService: MemberService;

  static get styles(): CSSResult[] {
    return [userDetailStyle, coreStyle];
  }

  firstUpdated(): void {
    this.getMemberByEmail();
  }

  openRFIDModal(): void {
    showComponent('#rfid-modal', this.shadowRoot);
  }

  getMemberByEmail(): void {
    this.memberService.getUsersMemberInfo().subscribe({
      next: (response: Member) => {
        this.memberUser = response;
        this.requestUpdate();
      },
    });
  }

  displayMemberResources(memberResources: MemberResource[]): TemplateResult {
    if (memberResources) {
      return html`
        <ol>
          ${memberResources?.map((x: MemberResource) => {
            return html` <li>${x.name}</li> `;
          })}
        </ol>
        <div class="lenny-face">
          ${memberResources?.length <= 3 ? '| (• ◡•)|' : ''}
        </div>
      `;
    } else {
      return html`
        <div>
          <span>You do not have any resources. </span>
          <div class="lenny-face">╥﹏╥</div>
        </div>
      `;
    }
  }

  refreshMember(): void {
    this.getMemberByEmail();
  }

  render(): TemplateResult {
    return html`
      <div class="user-profile">
        <paper-card>
          <div class="card-content">
            <h2>${this.memberUser?.name}</h2>
            <dl>
              <dt>Email</dt>
              <dd>${this.memberUser?.email}</dd>
              <dt>Status</dt>
              <dd>${displayMemberStatus(this.memberUser?.memberLevel)}</dd>
              <dt>Subscription ID</dt>
              <dd>${this.memberUser?.subscriptionID}</dd>
              <dt>RFID</dt>
              <dd>
                ${this.memberUser?.rfid !== 'notset'
                  ? this.memberUser?.rfid
                  : 'Not set'}
              </dd>
            </dl>
          </div>
          <div class="card-actions text-center">
            <mwc-button label="Assign RFID" @click=${this.openRFIDModal}>
            </mwc-button>
          </div>
        </paper-card>
        <paper-card>
          <div class="card-content">
            <h2>Your resources</h2>
            ${this.displayMemberResources(this.memberUser?.resources)}
          </div>
        </paper-card>
      </div>
      <rfid-modal
        id="rfid-modal"
        .email=${this.email}
        .showNewMemberOption=${false}
        .isThisSelf=${true}
        @updated=${this.refreshMember}
      >
      </rfid-modal>
    `;
  }
}
