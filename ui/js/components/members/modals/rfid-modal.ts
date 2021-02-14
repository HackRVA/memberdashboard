import { RFIDModalData } from "./../types/member.interface";
import { html, TemplateResult } from "lit-element";

export const rfidModal = (modalData: RFIDModalData): TemplateResult => {
  return html`
    <mwc-dialog id="assignRFIDModal">
      <div>Assign RFID</div>
      <mwc-textfield
        @change=${modalData.handleEmailChange}
        label="email"
        helper="member's email"
        value=${modalData.email}
      ></mwc-textfield>
      <mwc-textfield
        @change=${modalData.handleRFIDChange}
        label="RFID"
        helper="RFID"
        value=${modalData.rfid}
      ></mwc-textfield>
      <mwc-button
        slot="primaryAction"
        dialogAction="ok"
        @click=${modalData.handleSubmitForAssigningMemberToRFID}
      >
        Submit
      </mwc-button>
      <mwc-button
        slot="secondaryAction"
        dialogAction="cancel"
        @click=${modalData.emptyFormValuesOnClosed}
      >
        Cancel
      </mwc-button>
    </mwc-dialog>
  `;
};
