// lit element
import {
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// material
import { Dialog } from "@material/mwc-dialog";
import { Select } from "@material/mwc-select";
import { TextField } from "@material/mwc-textfield";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/select";

// membership
import { ResourceService } from "./../../../service/resource.service";
import { isEmpty } from "./../../../function";
import { MemberResource } from "../types";
import { RemoveMemberResourceRequest } from "../../resources/types";

@customElement("remove-member-from-resource-modal")
export class RemoveMemberFromResourceModal extends LitElement {
  @property({ type: String })
  email: string = "";

  @property({ type: Array })
  memberResources: Array<MemberResource> = [];

  resourceService: ResourceService = new ResourceService();

  removeResourceFromMemberModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  memberResourceSelectTemplate: Select;

  firstUpdated(): void {
    this.removeResourceFromMemberModalTemplate = this.shadowRoot.querySelector(
      "mwc-dialog"
    );
    this.emailFieldTemplate = this.shadowRoot.querySelector("mwc-textfield");
    this.memberResourceSelectTemplate = this.shadowRoot.querySelector(
      "mwc-select"
    );
  }

  show(): void {
    this.removeResourceFromMemberModalTemplate.show();
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.tryToRemoveMemberFromResource();
      this.emptyFormField();
      this.removeResourceFromMemberModalTemplate.close();
    } else {
      console.error("hrrrmm");
    }
  }

  tryToRemoveMemberFromResource(): void {
    const request: RemoveMemberResourceRequest = {
      email: this.emailFieldTemplate.value,
      resourceID: this.memberResourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.removeMemberFromResource(request);
  }

  removeMemberFromResource(request: RemoveMemberResourceRequest): void {
    this.resourceService.removeMemberFromResource(request).subscribe({
      complete: () => {
        const updatedEvent = new CustomEvent("updated");
        this.dispatchEvent(updatedEvent);
      },
    });
  }

  handleClosed(event: CustomEvent): void {
    // temp hack to stop mwc-select from bubbling to mwc-dialog
    const tagName: string = (event.target as EventTarget & { tagName: string })
      .tagName;
    if (tagName === "MWC-SELECT") {
      return;
    } else {
      this.emptyFormField();
    }
  }

  emptyFormField(): void {
    this.memberResourceSelectTemplate.select(-1);
  }

  isValid(): boolean {
    return (
      !isEmpty(this.emailFieldTemplate.value) &&
      !isEmpty(this.memberResourceSelectTemplate.value)
    );
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Remove Resource" @closed=${this.handleClosed}>
        <mwc-textfield
          label="email"
          helper="Can't edit email"
          value=${this.email}
          readonly
        ></mwc-textfield>
        <mwc-select label="Resources">
          ${this.memberResources?.map((x: MemberResource) => {
            return html`
              <mwc-list-item value=${x.resourceID}> ${x.name} </mwc-list-item>
            `;
          })}
        </mwc-select>
        <mwc-button slot="primaryAction" @click=${this.handleSubmit}>
          Submit
        </mwc-button>
        <mwc-button slot="secondaryAction" dialogAction="cancel">
          Cancel
        </mwc-button>
      </mwc-dialog>
    `;
  }
}
