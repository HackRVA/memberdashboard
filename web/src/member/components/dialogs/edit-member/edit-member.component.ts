// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Dialog } from '@material/mwc-dialog';
import { TextField } from '@material/mwc-textfield';

// memberdashboard
import '../../../../shared/components/toast-msg';
import { editMemberStyle } from './edit-member.style';
import {
  MemberService,
  MemberManagerService,
} from '../../../services/member.service';
import { Inject } from '../../../../shared/di/inject';
import { UpdateMemberRequest } from '../../../types/api/update-member-request';
import { displayToast } from '../../../../shared/components/abstract-toast';

@customElement('edit-member-form')
export class EditMemberForm extends LitElement {
  @property({ type: String })
  email: string = '';

  @property({ type: String })
  currentFullName: string = '';

  @property({ type: String })
  currentSubscriptionID: string = '';

  @property({ type: Function })
  closeHandler: () => void;

  @Inject('member')
  private memberService: MemberService;

  @Inject('member-manager')
  private memberManagerService: MemberManagerService;

  editMemberModalTemplate: Dialog;

  fullNameTemplate: TextField;

  subscriptionIDTemplate: TextField;

  static get styles(): CSSResult[] {
    return [editMemberStyle];
  }

  firstUpdated(): void {
    this.editMemberModalTemplate = this.shadowRoot.querySelector('mwc-dialog');
    this.fullNameTemplate = this.shadowRoot.querySelector('#name');
    this.subscriptionIDTemplate =
      this.shadowRoot.querySelector('#subscription-id');
  }

  updated(): void {
    this.fullNameTemplate.value = this.currentFullName;
    this.subscriptionIDTemplate.value = this.currentSubscriptionID;
  }

  private emptyFormField(): void {
    // fields are readonly
    this.fullNameTemplate.value = '';
  }

  private isValid(): boolean {
    return (
      this.fullNameTemplate.validity.valid &&
      !!this.fullNameTemplate.value.length
    );
  }

  private tryToUpdateMember(): void {
    const request: UpdateMemberRequest = {
      fullName: this.fullNameTemplate.value,
      subscriptionID: this.subscriptionIDTemplate.value,
    };

    this.updateMember(this.email, request);
  }

  private updateMember(email: string, request: UpdateMemberRequest): void {
    this.memberService.updateMemberByEmail(email, request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
        this.memberManagerService.getMembers();
        displayToast('success');
      },
      error: () => {
        displayToast('Unable to update member');
      },
    });
  }

  private fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
  }

  private handleSubmit(): void {
    if (this.isValid()) {
      this.tryToUpdateMember();
      this.emptyFormField();
    } else {
      displayToast('Hrmmm, are you sure everything in the form is correct?');
    }
  }

  render(): TemplateResult {
    return html`
      <mwc-textfield
        required
        type="text"
        outlined
        label="Full Name"
        helper="Full Name"
        id="name"
      ></mwc-textfield>
      <mwc-textfield
        required
        type="text"
        outlined
        label="Subscription ID"
        helper="Subscription ID"
        id="subscription-id"
      ></mwc-textfield>
      <mwc-button
        slot="primaryAction"
        @click=${() => {
          this.handleSubmit();
          this.closeHandler();
        }}
      >
        Submit
      </mwc-button>
      <mwc-button slot="secondaryAction" @click=${this.closeHandler}>
        Cancel
      </mwc-button>
    `;
  }
}
