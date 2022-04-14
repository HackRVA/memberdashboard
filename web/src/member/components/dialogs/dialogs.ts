import { html } from 'lit';

import {
    Member,
} from '../../types/api/member-response';

export default (shadowRoot) => {
    return {
        editMember: (member: Member): void => {
            const el = document.createElement('abstract-dialog');
            el.heading = "edit member";
            el.dialogLayout = html`
                <edit-member-form
                    id="edit-member-form"
                    .email=${member.email}
                    .currentFullName=${member.name}
                    .currentSubscriptionID=${member.subscriptionID}
                    .closeHandler=${() => {
                        el.hide();
                }}
                >
                </edit-member-form>`;
            shadowRoot.appendChild(el);
        },
        editRFID: (isNewMember) => {
            return (member: Member): void => {
            const el = document.createElement('abstract-dialog');
            el.heading = "edit RFID";
            el.dialogLayout = html`
                <rfid-form
                    id="edit-member-form"
                    .email=${member.email}
                    .RFID=${member.rfid}
                    .isNewMember=${isNewMember}
                    .showNewMemberOption=${false}
                    .closeHandler=${() => {
                        el.hide();
                }}
                >
                </rfid-form>`;

            shadowRoot.appendChild(el);
        }},
        addMemberToResource: (member: Member): void => {
            const el = document.createElement('abstract-dialog');
            el.heading = "add member to resource";
            el.dialogLayout = html`
                <add-resource-form
                    .emails=${[member.email]}
                    .closeHandler=${() => {
                        el.hide();
                }}
                >
                </add-resource-form>`;

            shadowRoot.appendChild(el);
        },
        removeMemberFromResource: (member: Member): void => {
            const el = document.createElement('abstract-dialog');
            el.heading = "remove member from resource";
            el.dialogLayout = html`
                <remove-resource-form
                    .email=${member.email}
                    .memberResources=${member.resources}
                    .closeHandler=${() => {
                        el.hide();
                }}
                >
                </remove-resource>`;

            shadowRoot.appendChild(el);
        }
    }
}
