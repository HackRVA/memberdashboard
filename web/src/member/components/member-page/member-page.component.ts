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
import { MemberManagerService } from '../../services/member.service';
import { Inject } from '../../../shared/di';

@customElement('member-page')
export class MemberPage extends LitElement {
  @property()
  members: Member[];

  @Inject('member-manager')
  private memberManagerService: MemberManagerService;

  static get styles(): CSSResult[] {
    return [coreStyle, memberPageStyle];
  }

  async firstUpdated(): Promise<void> {
    this.memberManagerService.registerListener(this.updateGrid);
    this.refreshMembers();
  }

  updateGrid = (): void => {
    this.members = this.memberManagerService.filteredMembers;
  };

  refreshMembers = (): void => {
    this.members = null;
    this.memberManagerService.getMembers();
  };

  render(): TemplateResult {
    return html`
      <loading-content .finishedLoading=${!!this.members}>
        <member-manager-menu></member-manager-menu>
        <member-search .updateGrid=${this.updateGrid}></member-search>
        <member-grid .members=${this.members}></member-grid>
      </loading-content>
    `;
  }
}
