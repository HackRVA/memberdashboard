import { RemoveMemberResourceModalData, MemberResource } from "./../types";
import { html, TemplateResult } from "lit-element";

export const removeMemberResourceModal = (
  modalData: RemoveMemberResourceModalData
): TemplateResult => {
  return html`
    <mwc-dialog id="removeMemberResourceModal">
      <div>Remove Resource</div>
      <mwc-textfield
        label="email"
        helper="Can't edit email"
        readonly
        value=${modalData.email}
      ></mwc-textfield>
      <mwc-select label="Resources" @change=${modalData.handleResourceChange}>
        ${modalData.memberResources.map((x: MemberResource) => {
          return html`
            <mwc-list-item value=${x.resourceID}> ${x.name} </mwc-list-item>
          `;
        })}
      </mwc-select>
      <mwc-button
        slot="primaryAction"
        dialogAction="ok"
        @click=${modalData.handleSubmitRemoveMemberResource}
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
