// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// material
import { TextField } from '@material/mwc-textfield/mwc-textfield';
import { Dialog } from '@material/mwc-dialog';

// memberdashboard
import '../toast-msg';
import { MemberService } from '../../../member/services/member.service';
import { ToastMessage } from '../../types/custom/toast-msg';
import { coreStyle } from '../../styles/core.style';
import { AssignRFIDRequest } from '../../../member/types/api/assign-rfid-request';
import { CreateMemberRequest } from '../../../member/types/api/create-member-request';
import { showComponent } from '../../functions';
import { rfidModalStyle } from './rfid-modal.style';
import { Inject } from '../../di';
import { IPopup } from './../../types/custom/ipop-up';

@customElement('rfid-modal')
export class RFIDModal extends LitElement implements IPopup {
  @property({ type: String })
  email: string = '';

  @property({ type: Boolean })
  showNewMemberOption: boolean = false;

  @property({ type: Boolean })
  isThisSelf: boolean = false;

  @Inject('member')
  private memberService: MemberService;

  rfidModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  rfidFieldTemplate: TextField;

  toastMsg: ToastMessage;

  isNewMember: Boolean;

  static get styles(): CSSResult[] {
    return [rfidModalStyle, coreStyle];
  }

  firstUpdated(): void {
    this.rfidModalTemplate = this.shadowRoot?.querySelector('mwc-dialog');
    this.emailFieldTemplate = this.shadowRoot?.querySelector('#email');
    this.rfidFieldTemplate = this.shadowRoot?.querySelector('#rfid');
  }

  updated(): void {
    this.emailFieldTemplate.value = this.email;

    if (this.isThisSelf) {
      this.emailFieldTemplate.disabled = true;
    }
  }

  public show(): void {
    this.rfidModalTemplate?.show();
  }

  private tryToAssigningMemberToRFID(): void {
    const request: AssignRFIDRequest | CreateMemberRequest = {
      email: this.emailFieldTemplate.value.trim(),
      rfid: this.rfidFieldTemplate.value.trim(),
    };

    if (this.isNewMember) {
      this.assignNewMemberToRFID(request);
      return;
    }
    this.assignMemberToRFID(request);
  }

  private tryToAssigningSelfToRFID(): void {
    const request: AssignRFIDRequest = {
      email: this.emailFieldTemplate.value.trim(),
      rfid: this.rfidFieldTemplate.value.trim(),
    };

    this.assignRFIDToSelf(request);
  }

  private assignNewMemberToRFID(request: CreateMemberRequest): void {
    this.memberService.assignNewMemberRFID(request).subscribe({
      complete: () => {
        this.displayToastMsg('Success');
        this.fireUpdatedEvent();
        this.emptyFormField();
        this.rfidModalTemplate.close();
      },
      error: () => {
        this.displayToastMsg('Hrmmm, something went wrong? :3');
      },
    });
  }

  private assignMemberToRFID(request: AssignRFIDRequest): void {
    this.memberService.assignRFID(request).subscribe({
      complete: () => {
        this.displayToastMsg('Success');
        this.fireUpdatedEvent();
        this.emptyFormField();
        this.rfidModalTemplate.close();
      },
      error: () => {
        this.displayToastMsg('Hrmmm, are you sure this is a member? :3');
      },
    });
  }

  private assignRFIDToSelf(request: AssignRFIDRequest): void {
    this.memberService.assignRFIDToSelf(request).subscribe({
      complete: () => {
        this.displayToastMsg('Success');
        this.fireUpdatedEvent();
        this.rfidModalTemplate.close();
      },
      error: () => {
        this.displayToastMsg('Hrmmm, are you sure this is a member? :3');
      },
    });
  }

  private handleNewMember(): void {
    this.isNewMember = !this.isNewMember;
  }

  private handleSubmit(): void {
    if (this.isValid()) {
      if (this.isThisSelf) {
        this.tryToAssigningSelfToRFID();
      } else {
        this.tryToAssigningMemberToRFID();
      }
    } else {
      this.displayToastMsg(
        'Hrmmm, are you sure everything in the form is correct?'
      );
    }
  }

  private fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
  }

  private emptyFormField(): void {
    // fields are readonly
    if (!this.email) {
      this.emailFieldTemplate.value = '';
    }
    this.rfidFieldTemplate.value = '';
  }

  private isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.rfidFieldTemplate.validity.valid
    );
  }

  private handleClosed(): void {
    this.emptyFormField();
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  private displayNewMemberCheckBox(): TemplateResult {
    if (!this.showNewMemberOption) return html``;

    return html`
      <mwc-formfield label="New member" class="new-member">
        <mwc-checkbox @change=${this.handleNewMember}></mwc-checkbox>
      </mwc-formfield>
    `;
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Assign RFID" @closed=${this.handleClosed}>
        <mwc-textfield
          required
          type="email"
          label="email"
          helper="email"
          id="email"
        ></mwc-textfield>
        <mwc-textfield
          required
          label="RFID"
          helper="RFID"
          id="rfid"
          type="number"
        ></mwc-textfield>
        ${this.displayNewMemberCheckBox()}
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
