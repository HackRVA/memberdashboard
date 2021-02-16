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
import { Select } from "@material/mwc-select";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/select";

// membership
import { ResourceService } from "../../../service/resource.service";
import {
  AddMemberResourceRequest,
  ResourceResponse,
} from "../../resources/types";
import { isEmpty } from "../../../function";

@customElement("add-member-to-resource-modal")
export class AddMemberToResourceModal extends LitElement {
  @property({ type: String })
  email: string = "";

  resourceService: ResourceService = new ResourceService();
  resources: ResourceResponse[] = [];

  addResourceToMemberModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  resourceSelectTemplate: Select;

  firstUpdated(): void {
    this.addResourceToMemberModalTemplate = this.shadowRoot.querySelector(
      "mwc-dialog"
    );
    this.emailFieldTemplate = this.shadowRoot.querySelector("mwc-textfield");
    this.resourceSelectTemplate = this.shadowRoot.querySelector("mwc-select");
    this.getResources();
  }

  show(): void {
    this.addResourceToMemberModalTemplate?.show();
  }

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any })?.error) {
          console.error("some error getting resources");
        } else {
          this.resources = result as ResourceResponse[];
          this.requestUpdate();
        }
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

  handleSubmit(): void {
    if (this.isValid()) {
      this.tryToAddMemberToResource();
      this.emptyFormField();
      this.addResourceToMemberModalTemplate.close();
    } else {
      console.error("hrrmmmm");
    }
  }

  tryToAddMemberToResource(): void {
    const request: AddMemberResourceRequest = {
      email: this.emailFieldTemplate.value,
      resourceID: this.resourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.addMemberToResource(request);
  }

  addMemberToResource(request: AddMemberResourceRequest): void {
    this.resourceService.addMemberToResource(request).subscribe({
      complete: () => {
        const updatedEvent = new CustomEvent("updated");
        this.dispatchEvent(updatedEvent);
      },
    });
  }

  emptyFormField(): void {
    this.resourceSelectTemplate.select(-1);
  }

  isValid(): boolean {
    return (
      !isEmpty(this.emailFieldTemplate.value) &&
      !isEmpty(this.resourceSelectTemplate.value)
    );
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Add Resource" @closed=${this.handleClosed}>
        <mwc-textfield
          label="email"
          helper="Can't edit email"
          value=${this.email}
          readonly
        ></mwc-textfield>
        <mwc-select label="Resources">
          ${this.resources.map((x: ResourceResponse) => {
            return html`
              <mwc-list-item value=${x.id}> ${x.name} </mwc-list-item>
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
