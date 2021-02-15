// lit element
import { customElement, html, LitElement, TemplateResult } from "lit-element";

// material
import { TextField } from "@material/mwc-textfield/mwc-textfield";
import { Dialog } from "@material/mwc-dialog";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-select";

// membership
import { MemberService } from "../../service/member.service";
import { AssignRFIDRequest } from "../members/types";
import { defaultSnackbar } from "./default-snackbar";
import { showComponent } from "../../function";

@customElement("rfid-modal")
export class RFIDModal extends LitElement {
  memberService: MemberService = new MemberService();

  rfidModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  rfidFieldTemplate: TextField;

  firstUpdated(): void {
    this.rfidModalTemplate = this.shadowRoot?.querySelector("mwc-dialog");
    this.emailFieldTemplate = this.shadowRoot?.querySelector("#email");
    this.rfidFieldTemplate = this.shadowRoot?.querySelector("#rfid");
  }

  show(): void {
    this.rfidModalTemplate?.show();
  }

  tryToAssigningMemberToRFID(): void {
    const request: AssignRFIDRequest = {
      email: this.emailFieldTemplate.value.trim(),
      rfid: this.rfidFieldTemplate.value,
    };

    this.assignMemberToRFID(request);
  }

  assignMemberToRFID(request: AssignRFIDRequest): void {
    this.memberService.assignRFID(request).subscribe({
      complete: () => {
        this.displaySuccessMessage();
        this.requestUpdate();
      },
    });
  }

  displaySuccessMessage(): void {
    showComponent("#success", this.shadowRoot);
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.tryToAssigningMemberToRFID();
      this.emptyFormField();
      this.rfidModalTemplate.close();
    } else {
      console.error("Hrmmmm");
    }
  }

  emptyFormField(): void {
    // fields are readonly
    this.emailFieldTemplate.value = "";
    this.rfidFieldTemplate.value = "";
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.rfidFieldTemplate.validity.valid
    );
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Assign RFID">
        <mwc-textfield
          required
          label="email"
          helper="email"
          id="email"
        ></mwc-textfield>
        <mwc-textfield
          required
          label="RFID"
          helper="RFID"
          id="rfid"
        ></mwc-textfield>
        <mwc-button slot="primaryAction" @click=${this.handleSubmit}>
          Submit
        </mwc-button>
        <mwc-button slot="secondaryAction" dialogAction="cancel">
          Cancel
        </mwc-button>
      </mwc-dialog>
      ${defaultSnackbar("success", "success")}
    `;
  }
}
