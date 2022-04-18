// lit element
import { customElement, property, state } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// vaadin
import { TextFieldValueChangedEvent } from '@vaadin/text-field';

// memberdashboard
import { Member } from '../../types/api/member-response';
import { displayMemberStatus } from '../../functions';
import { MemberManagerService } from '../../services/member.service';
import { Inject } from '../../../shared/di/inject';

@customElement('member-search')
export class MemberSearch extends LitElement {
    @Inject('member-manager')
    private memberManagerService: MemberManagerService;

    @property({ type: Function })
    updateGrid: Function;

    render(): TemplateResult {
        return html`
        <vaadin-text-field
            placeholder="Search"
            style="width: 50%;"
            @value-changed="${(e: TextFieldValueChangedEvent) => {
                const searchTerm = ((e.detail.value as string) || '').trim();
                const matchesTerm = (value: string) => {
                    return value.toLowerCase().indexOf(searchTerm.toLowerCase()) >= 0;
                };

                const members = this.memberManagerService?.showActive ? this.memberManagerService?.activeMembers : this.memberManagerService?.inactiveMembers
                this.memberManagerService.filteredMembers = (members || [])
                    .filter(({ name, email, rfid, memberLevel }) => {
                        return (
                            !searchTerm ||
                            matchesTerm(name) ||
                            matchesTerm(email) ||
                            matchesTerm(rfid) ||
                            matchesTerm(displayMemberStatus(memberLevel))
                        );
                    });

                this.updateGrid();
            }}"
        >
            <vaadin-icon slot="prefix" icon="vaadin:search"></vaadin-icon>
        </vaadin-text-field>
        `;
    }
}
