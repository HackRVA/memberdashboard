// lit element
import { customElement, property, state } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../../../shared/components/rfid-form';
import '../../../shared/components/abstract-dialog';
import '../dialogs/edit-member';
import '../dialogs/remove-resource';
import '../dialogs/add-resource';
import { Member } from '../../types/api/member-response';
import { coreStyle } from '../../../shared/styles';
import { memberGridStyle } from './member-grid.style';
import { displayMemberStatus } from '../../functions';
import dialogs from '../dialogs/dialogs';

function actionBuilder(getLabel: Function, getAction: Function, renderer?: Function) {
    if (renderer) return renderer;

    return (root: HTMLElement, grid, model): void => {
        if (root.firstElementChild) return;
        const btn: any = document.createElement('vaadin-button')
        btn.addEventListener('click', () => {
            getAction(model.item)
        })
        btn.textContent = getLabel(model.item);
        root.appendChild(btn);
    }
}

@customElement('member-grid')
export class MemberGrid extends LitElement {
    @property({type: Array})
    members: Member[] = [];

    grid: any;

    static get styles(): CSSResult[] {
        return [memberGridStyle, coreStyle];
    }

    private actions: {
        text: string;
    }[] = [
            { text: 'Assign RFID' },
            { text: 'Edit Member' },
            { text: 'Add resource' },
            { text: 'Remove resource' },
        ];

    private resourceActions: {
        text: string;
    }[] = [
            { text: 'Add resource' },
            { text: 'Remove resource' },
        ];

    dialogs: any;

    async firstUpdated(): Promise<void> {
        this.grid = this.shadowRoot.querySelector('#member-grid')
        this.dialogs = dialogs(this.shadowRoot)
    }

    private menuEventHandler = (model: any) => {
        return (e: any) => {
            if (!this.dialogs) return;
            const actionHandlers = {
                'Assign RFID': this.dialogs.editRFID(false),
                'Edit Member': this.dialogs.editMember,
                'Add resource': this.dialogs.addMemberToResource,
                'Remove resource': this.dialogs.removeMemberFromResource,
            };

            actionHandlers[e.detail.value.text](model.item)
        };
    }

    private makeIcon(): HTMLElement {
        const item = document.createElement('vaadin-context-menu-item');
        item.textContent = '•••';
        item.setAttribute('aria-label', 'More options');
        return item;
    }

    private actionsRenderer = (root: HTMLElement, grid, model): void => {
        if (root.firstElementChild) return;

        const menuBar: any = document.createElement('vaadin-menu-bar');
        menuBar.items = [{ component: this.makeIcon(), children: this.actions, }];
        menuBar.addEventListener('item-selected', this.menuEventHandler(model));
        menuBar.setAttribute('theme', 'tertiary');
        root.appendChild(menuBar);
    }

    render(): TemplateResult {
        if (!this.dialogs) return;
        return html`
        <vaadin-grid
            id="member-grid" 
            .items="${this.members}" 
            theme="row-dividers" 
            column-reordering-allowed 
            style="height: 70vh;"
        >
            <vaadin-grid-column
                path="name"
                .renderer="${actionBuilder((member) => member.name, this.dialogs.editMember)}"
            ></vaadin-grid-column>
            <vaadin-grid-column auto-width flex-grow="0" path="email"></vaadin-grid-column>
            <vaadin-grid-column auto-width flex-grow="0" path="rfid"
                .renderer="${actionBuilder((member) => member.rfid, this.dialogs.editRFID(false))}"
            ></vaadin-grid-column>
            <vaadin-grid-column auto-width flex-grow="0" id="member-level" path="memberLevel"
                .renderer="${actionBuilder((member) => displayMemberStatus(member.memberLevel), () => console.log("show some info"))}"
            ></vaadin-grid-column>
            <vaadin-grid-column id="resources" path="resources" 
                .renderer="${(root: HTMLElement, grid, model): void => {
                    if (root.firstElementChild) return;
            
                    const label = document.createElement('vaadin-context-menu-item')
                    label.textContent = model.item.resourcesLabel || 'no resources';
            
                    const menuBar: any = document.createElement('vaadin-menu-bar');
                    menuBar.items = [{ component: label, children: this.resourceActions, }];
                    menuBar.addEventListener('item-selected', this.menuEventHandler(model))
                    menuBar.setAttribute('theme', 'tertiary');
                    root.appendChild(menuBar);
                }}"></vaadin-grid-column>
            <vaadin-grid-column auto-width flex-grow="0" path="actions" 
                .renderer="${this.actionsRenderer}"></vaadin-grid-column>
        </vaadin-grid>
        `
    }
}
