// lit element
import {
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// material
import { TextField } from "@material/mwc-textfield/mwc-textfield";
import { Dialog } from "@material/mwc-dialog";

// membership
import { MemberService } from "../../service";
import { AssignRFIDRequest, CreateMemberRequest } from "../members/types";
import { showComponent } from "../../function";
import { ToastMessage } from "./types";
import "./toast-msg";

@customElement("rfid-modal")
export class RFIDModal extends LitElement {
  @property({ type: String })
  email: string = "";
  @property({ type: Boolean })
  showNewMemberOption: boolean = false;

  memberService: MemberService = new MemberService();

  rfidModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  rfidFieldTemplate: TextField;

  toastMsg: ToastMessage;

  isNewMember: Boolean;

  firstUpdated(): void {
    this.rfidModalTemplate = this.shadowRoot?.querySelector("mwc-dialog");
    this.emailFieldTemplate = this.shadowRoot?.querySelector("#email");
    this.rfidFieldTemplate = this.shadowRoot?.querySelector("#rfid");
  }

  updated(): void {
    if (this.email) {
      this.emailFieldTemplate.value = this.email;
      this.emailFieldTemplate.disabled = true;
    }
  }

  show(): void {
    this.rfidModalTemplate?.show();
  }

  tryToAssigningMemberToRFID(): void {
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

  tryToAssigningSelfToRFID(): void {
    const request: AssignRFIDRequest = {
      email: this.emailFieldTemplate.value.trim(),
      rfid: this.rfidFieldTemplate.value.trim(),
    };

    this.assignRFIDToSelf(request);
  }

  assignNewMemberToRFID(request: CreateMemberRequest): void {
    this.memberService.assignNewMemberRFID(request).subscribe({
      complete: () => {
        this.displayToastMsg("Success");
        this.fireUpdatedEvent();
        this.emptyFormField();
        this.rfidModalTemplate.close();
      },
      error: () => {
        this.displayToastMsg("Hrmmm, something went wrong? :3");
      },
    });
  }

  assignMemberToRFID(request: AssignRFIDRequest): void {
    this.memberService.assignRFID(request).subscribe({
      complete: () => {
        this.displayToastMsg("Success");
        this.fireUpdatedEvent();
        this.emptyFormField();
        this.rfidModalTemplate.close();
      },
      error: () => {
        this.displayToastMsg("Hrmmm, are you sure this is a member? :3");
      },
    });
  }

  assignRFIDToSelf(request: AssignRFIDRequest): void {
    this.memberService.assignRFIDToSelf(request).subscribe({
      complete: () => {
        this.displayToastMsg("Success");
        this.fireUpdatedEvent();
        this.rfidModalTemplate.close();
      },
      error: () => {
        this.displayToastMsg("Hrmmm, are you sure this is a member? :3");
      },
    });
  }

  handleNewMember(): void {
    this.isNewMember = !this.isNewMember;
  }

  handleSubmit(): void {
    if (this.isValid()) {
      if (this.email) {
        this.tryToAssigningSelfToRFID();
      } else {
        this.tryToAssigningMemberToRFID();
      }
    } else {
      this.displayToastMsg(
        "Hrmmm, are you sure everything in the form is correct?"
      );
    }
  }

  fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent("updated");
    this.dispatchEvent(updatedEvent);
  }

  emptyFormField(): void {
    // fields are readonly
    if (!this.email) {
      this.emailFieldTemplate.value = "";
    }
    this.rfidFieldTemplate.value = "";
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.rfidFieldTemplate.validity.valid
    );
  }

  handleClosed(): void {
    this.emptyFormField();
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  displayNewMemberCheckBox(): TemplateResult {
    if (!this.showNewMemberOption) return html``;

    return html`
      <mwc-formfield label="New member">
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
