// lit-element
import { html } from 'lit';

// memberdashboard
import { Member } from '../../types/api/member-response';

export default (
  shadowRoot: ShadowRoot
): {
  editMember: (member: Member) => void;
  editRFID: (isNewMember: boolean) => void;
  addMemberToResource: (member: Member) => void;
  removeMemberFromResource: (member: Member) => void;
  refreshMemberStatus: (member: Member) => void;
} => {
  return {
    editMember: (member: Member): void => {
      const el = document.createElement('abstract-dialog');
      el.heading = 'edit member';
      el.dialogLayout = html` <edit-member-form
        id="edit-member-form"
        .email=${member.email}
        .currentFullName=${member.name}
        .currentSubscriptionID=${member.subscriptionID}
        .closeHandler=${(): void => {
          el.hide();
        }}
      >
      </edit-member-form>`;
      shadowRoot.appendChild(el);
    },
    editRFID: (isNewMember: boolean): ((member: Member) => void) => {
      return (member: Member): void => {
        const el = document.createElement('abstract-dialog');
        el.heading = 'Edit RFID';
        el.dialogLayout = html` <rfid-form
          id="edit-member-form"
          .email=${member.email}
          .RFID=${member.rfid}
          .isNewMember=${isNewMember}
          .showNewMemberOption=${false}
          .closeHandler=${(): void => {
            el.hide();
          }}
        >
        </rfid-form>`;

        shadowRoot.appendChild(el);
      };
    },
    addMemberToResource: (member: Member): void => {
      const el = document.createElement('abstract-dialog');
      el.heading = 'add member to resource';
      el.dialogLayout = html` <add-resource-form
        .emails=${[member.email]}
        .closeHandler=${(): void => {
          el.hide();
        }}
      >
      </add-resource-form>`;

      shadowRoot.appendChild(el);
    },
    removeMemberFromResource: (member: Member): void => {
      const el = document.createElement('abstract-dialog');
      el.heading = 'remove member from resource';
      el.dialogLayout = html`
                <remove-resource-form
                    .email=${member.email}
                    .memberResources=${member.resources}
                    .closeHandler=${(): void => {
                      el.hide();
                    }}
                >
                </remove-resource>`;

      shadowRoot.appendChild(el);
    },
    refreshMemberStatus: (member: Member): void => {
      const el = document.createElement('abstract-dialog');
      el.heading = 'member status';
      el.dialogLayout = html` <refresh-member-status-form
        .subscriptionID=${member.subscriptionID}
        .closeHandler=${(): void => {
          el.hide();
        }}
      >
      </refresh-member-status-form>`;

      shadowRoot.appendChild(el);
    },
  };
};
