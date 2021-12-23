// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Dialog } from '@material/mwc-dialog';
import { TextField } from '@material/mwc-textfield';

// memberdashboard
import { editMemberStyle } from './edit-member.style';

@customElement('edit-member')
export class AddMembersToResourceModal extends LitElement {
  @property({ type: String })
  email: string = '';

  @property({ type: String })
  currentFullName: string = '';

  editMemberModalTemplate: Dialog;
  fullNameTemplate: TextField;

  static get styles(): CSSResult[] {
    return [editMemberStyle];
  }

  firstUpdated(): void {
    this.editMemberModalTemplate = this.shadowRoot.querySelector('mwc-dialog');
    this.fullNameTemplate = this.shadowRoot.querySelector('#name');
  }

  updated(): void {
    this.fullNameTemplate.value = this.currentFullName;
  }

  show(): void {
    this.editMemberModalTemplate?.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Edit Member">
        <mwc-textfield
          required
          type="text"
          outlined
          label="Full Name"
          helper="Full Name"
          id="name"
        ></mwc-textfield>
        <mwc-button slot="primaryAction"> Submit </mwc-button>
        <mwc-button slot="secondaryAction" dialogAction="cancel">
          Cancel
        </mwc-button>
      </mwc-dialog>
      <toast-msg id="toast-msg"> </toast-msg>
    `;
  }
}
