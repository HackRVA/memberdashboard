// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// material
import { TextField } from '@material/mwc-textfield/mwc-textfield';

// memberdashboard
import '../toast-msg';
import { MemberService } from '../../../member/services/member.service';
import { coreStyle } from '../../styles/core.style';
import { AssignRFIDRequest } from '../../../member/types/api/assign-rfid-request';
import { CreateMemberRequest } from '../../../member/types/api/create-member-request';
import { displayToast } from '../../components/abstract-toast';
import { rfidFormStyle } from './rfid-form.style';
import { Inject } from '../../di';

@customElement('rfid-form')
export class RFIDModal extends LitElement {
  @property({ type: String })
  email: string = '';

  @property({ type: String })
  RFID: string = '';

  @property({ type: Boolean })
  showNewMemberOption: boolean = false;

  @property({ type: Boolean })
  isThisSelf: boolean = false;

  @property({ type: Function })
  closeHandler: () => void;

  @property({ type: Boolean })
  isNewMember: boolean;

  @Inject('member')
  private memberService: MemberService;

  emailFieldTemplate: TextField;
  rfidFieldTemplate: TextField;

  static get styles(): CSSResult[] {
    return [rfidFormStyle, coreStyle];
  }

  firstUpdated(): void {
    this.rfidFieldTemplate = this.shadowRoot?.querySelector('#rfid');
    this.emailFieldTemplate = this.shadowRoot?.querySelector('#email');
  }

  updated(): void {
    this.emailFieldTemplate.value = this.email || '';
    this.rfidFieldTemplate.value = this.RFID || '';

    if (this.isThisSelf) {
      this.emailFieldTemplate.disabled = true;
    }
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
        this.fireUpdatedEvent();
        this.emptyFormField();
        // this.memberService.getMembers();
        displayToast('Success');
      },
      error: () => {
        displayToast('Hrmmm, something went wrong? :3');
      },
    });
  }

  private assignMemberToRFID(request: AssignRFIDRequest): void {
    this.memberService.assignRFID(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
        this.emptyFormField();
        // this.memberService.getMembers();
        displayToast('Success');
      },
      error: () => {
        displayToast('Hrmmm, are you sure this is a member? :3');
      },
    });
  }

  private assignRFIDToSelf(request: AssignRFIDRequest): void {
    this.memberService.assignRFIDToSelf(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
        this.closeHandler();
        // this.memberService.getMembers();
        displayToast('Success');
      },
      error: () => {
        displayToast('Hrmmm, are you sure this is a member? :3');
      },
    });
  }

  private handleSubmit(): void {
    if (this.isValid()) {
      if (this.isThisSelf) {
        this.tryToAssigningSelfToRFID();
      } else {
        this.tryToAssigningMemberToRFID();
      }
      this.closeHandler();
    } else {
      displayToast('Hrmmm, are you sure everything in the form is correct?');
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

  private newMemberDisclaimer() {
    if (!this.isNewMember) return;
    return html`<span>
        This is an option if the new member's info hasn't been sent from paypal
        yet. Check the table before trying to add them manually. note: if you
        add the member with this method, their email has to match what they use
        in paypal.
      </span>
      <hr /> `;
  }

  render(): TemplateResult {
    return html`
      ${this.newMemberDisclaimer()}
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
      <mwc-button slot="primaryAction" @click=${this.handleSubmit}>
        Submit
      </mwc-button>
      <mwc-button
        slot="secondaryAction"
        dialogAction="cancel"
        @click=${this.closeHandler}
      >
        Cancel
      </mwc-button>
    `;
  }
}
