import { MemberService } from "./../../../service/member.service";
import { html, TemplateResult } from "lit-element";
import { ResourceService } from "../../../service/resource.service";

export const AddMemberResourceModal = (
  modalData: MemberService.AddMemberResourceModalData
): TemplateResult => {
  return html`
    <mwc-dialog id="addMemberResourceModal">
      <div>Add Resource</div>
      <mwc-textfield
        label="email"
        helper="Can't edit email"
        readonly
        value=${modalData.email}
      ></mwc-textfield>
      <mwc-select label="Resources" @change=${modalData.handleResourceChange}>
        ${modalData.resources.map((x: ResourceService.ResourceResponse) => {
          return html`
            <mwc-list-item value=${x.id}> ${x.name} </mwc-list-item>
          `;
        })}
      </mwc-select>
      <mwc-button
        slot="primaryAction"
        dialogAction="ok"
        @click=${modalData.handleSubmitAddMemberResource}
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
