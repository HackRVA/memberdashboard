// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import { memberManagerMenuStyle } from './member-manager-menu.style';
import { MemberManagerService } from '../../services/member.service';
import { Inject } from '../../../shared/di';
import dialogs from '../dialogs/dialogs';
import { coreStyle } from '../../../shared/styles';

@customElement('member-manager-menu')
export class MemberManagerMenu extends LitElement {
    @Inject('member-manager')
    private memberManagerService: MemberManagerService;

    static get styles(): CSSResult[] {
        return [coreStyle, memberManagerMenuStyle];
    }

    @property({ type: Function })
    refreshMembers: Function;

    render(): TemplateResult {
        return html`
        <div class="menuContainer">
            <vaadin-tabs>
                <vaadin-tab @click="${() => {
                this.memberManagerService.showActive = true;
            }}"
                >
                    Active
                </vaadin-tab>
                <vaadin-tab @click="${() => {
                this.memberManagerService.showActive = false;
            }}"
                >
                    Inactive
                </vaadin-tab>
            </vaadin-tabs>
            <vaadin-button @click=${dialogs(this.shadowRoot).editRFID(true)}>
                New Member
            </vaadin-button>
        </div>
        `;
    }
}
