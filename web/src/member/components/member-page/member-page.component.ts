// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../member-grid';
import '../member-search';
import '../member-manager-menu';
import '../../../shared/components/md-card';
import '../../../shared/components/loading-content';
import { memberPageStyle } from './member-page.style';
import { Member } from '../../types/api/member-response';
import { coreStyle } from '../../../shared/styles';
import { MemberService } from '../../services/member.service';
import { Inject } from '../../../shared/di';

@customElement('member-page')
export class MemberPage extends LitElement {
  @Inject('member')
  private memberService: MemberService;

  page = 0;
  count = 10;
  active = true;
  search = '';

  members: Member[];

  static get styles(): CSSResult[] {
    return [coreStyle, memberPageStyle];
  }

  firstUpdated(): void {
    this.refreshMembers();
  }

  async refreshMembers(): Promise<void> {
    this.members = await this.memberService
      .getMembers(this.page, this.count, this.active, this.search)
      .toPromise();
    this.requestUpdate();
  }

  incrementPage(): void {
    this.page++;
    this.refreshMembers();
  }
  decrementPage(): void {
    if (this.page == 0) {
      return;
    }
    this.page--;
    this.refreshMembers();
  }

  memberSeach = (searchTerm: string): Promise<void> => {
    return (async (): Promise<void> => {
      if (searchTerm.length == 0) {
        this.refreshMembers();
        return;
      }
      this.members = await this.memberService
        .getMembers(0, 0, this.active, searchTerm)
        .toPromise();
      this.requestUpdate();
    })();
  };

  setActive = (active: boolean): Promise<void> => {
    return (async (): Promise<void> => {
      this.active = active;
      this.refreshMembers();
    })();
  };

  render(): TemplateResult {
    return html`
      <loading-content .finishedLoading=${!!this.members}>
        <member-manager-menu .setActive=${this.setActive}></member-manager-menu>
        <member-search .memberSearch=${this.memberSeach}></member-search>
        <member-grid .members=${this.members}></member-grid>
        <button @click=${this.decrementPage}><</button>
        <button @click=${this.incrementPage}>></button>
        <span> ${this.page} </span>
      </loading-content>
    `;
  }
}
