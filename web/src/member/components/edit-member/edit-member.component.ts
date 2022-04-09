// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Dialog } from '@material/mwc-dialog';
import { TextField } from '@material/mwc-textfield';

// memberdashboard
import '../../../shared/components/toast-msg';
import { editMemberStyle } from './edit-member.style';
import { MemberService } from '../../services/member.service';
import { Inject } from '../../../shared/di/inject';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { showComponent } from '../../../shared/functions';
import { UpdateMemberRequest } from '../../types/api/update-member-request';
import { IPopup } from '../../../shared/types/custom/ipop-up';

@customElement('edit-member')
export class EditMemberModal extends LitElement implements IPopup {
  @property({ type: String })
  email: string = '';

  @property({ type: String })
  currentFullName: string = '';

  @property({ type: String })
  currentSubscriptionID: string = '';

  @Inject('member')
  private memberService: MemberService;

  toastMsg: ToastMessage;

  editMemberModalTemplate: Dialog;

  fullNameTemplate: TextField;

  subscriptionIDTemplate: TextField;

  static get styles(): CSSResult[] {
    return [editMemberStyle];
  }

  firstUpdated(): void {
    this.editMemberModalTemplate = this.shadowRoot.querySelector('mwc-dialog');
    this.fullNameTemplate = this.shadowRoot.querySelector('#name');
    this.subscriptionIDTemplate = this.shadowRoot.querySelector('#subscription-id');
  }

  updated(): void {
    this.fullNameTemplate.value = this.currentFullName;
    this.subscriptionIDTemplate.value = this.currentSubscriptionID;
  }

  private handleClosed(): void {
    this.emptyFormField();
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
      },
      error: () => {
        this.displayToastMsg('Unable to update member');
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
      this.editMemberModalTemplate.close();
    } else {
      this.displayToastMsg(
        'Hrmmm, are you sure everything in the form is correct?'
      );
    }
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  public show(): void {
    this.editMemberModalTemplate?.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Edit Member" @closed=${this.handleClosed}>
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
        <mwc-button slot="primaryAction" @click=${this.handleSubmit}>
          Submit
        </mwc-button>
        <mwc-button slot="secondaryAction" dialogAction="cancel">
          Cancel
        </mwc-button>
      </mwc-dialog>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
