// lit element
import { customElement } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../member-manager';
import '../../../shared/components/md-card';
import '../../../shared/components/loading-content';
import { coreStyle } from '../../../shared/styles';
import { MemberService } from '../../services/member.service';
import { MemberResponse } from '../../types/api/member-response';
import { MemberLevel } from '../../types/custom/member-level';
import { Inject } from '../../../shared/di';

@customElement('member-page')
export class MemberPage extends LitElement {
  members: MemberResponse[];
  memberCount: number = 0;
  totalMemberCount: number = 0;
  finishedLoading: boolean = false;

  @Inject('member')
  private memberService: MemberService;

  static get styles(): CSSResult[] {
    return [coreStyle];
  }

  firstUpdated(): void {
    this.getMembers();
  }

  getMembers(): void {
    this.memberService.getMembers().subscribe({
      next: (result: MemberResponse[]) => {
        this.finishedLoading = true;
        this.members = result;
        this.totalMemberCount = result.length;
        this.memberCount = this.getActiveMembers().length;
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to get members');
      },
    });
  }

  getActiveMembers(): MemberResponse[] {
    return this.members.filter(
      (x: MemberResponse) => x.memberLevel !== MemberLevel.inactive
    );
  }

  render(): TemplateResult {
    return html`
      <md-card>
        <loading-content .finishedLoading=${this.finishedLoading}>
          <member-manager
            .members=${this.members}
            .memberCount=${this.memberCount}
            .totalMemberCount=${this.totalMemberCount}
          ></member-manager>
        </loading-content>
      </md-card>
    `;
  }
}