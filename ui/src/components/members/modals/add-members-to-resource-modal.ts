// lit element
import {
  CSSResult,
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// material
import { Dialog } from "@material/mwc-dialog";
import { Select } from "@material/mwc-select";

// memberdashboard
import { ResourceService } from "../../../service";
import {
  BulkAddMembersToResourceRequest,
  ResourceResponse,
} from "../../resources/types";
import { isEmpty, showComponent } from "../../../function";
import { ToastMessage } from "../../shared/types";
import { addMembersToResourceModalStyles } from "../styles/add-members-to-resource-modal-styles";

@customElement("add-members-to-resource-modal")
export class AddMembersToResourceModal extends LitElement {
  @property({ type: Array })
  emails: string[] = [];

  addResourceToMembersModalTemplate: Dialog;
  resourceSelectTemplate: Select;

  resourceService: ResourceService = new ResourceService();
  resources: ResourceResponse[] = [];
  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [addMembersToResourceModalStyles];
  }

  firstUpdated(): void {
    this.addResourceToMembersModalTemplate =
      this.shadowRoot.querySelector("mwc-dialog");
    this.resourceSelectTemplate = this.shadowRoot.querySelector("mwc-select");

    this.getResources();
  }

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result: ResourceResponse[]) => {
        this.resources = result;
        this.requestUpdate();
      },
      error: () => {
        console.error("unable to get resources");
      },
    });
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  tryToAddMembersToResource(): void {
    const request: BulkAddMembersToResourceRequest = {
      emails: this.emails,
      resourceID: this.resourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.bulkAddMembersToResource(request);
  }

  bulkAddMembersToResource(request: BulkAddMembersToResourceRequest): void {
    this.resourceService.bulkAddMembersToResource(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
      },
    });
  }

  fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent("updated");
    this.dispatchEvent(updatedEvent);
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.tryToAddMembersToResource();
      this.emptyFormField();
      this.addResourceToMembersModalTemplate.close();
    } else {
      this.displayToastMsg("Hrmmmmm");
    }
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

  isValid(): boolean {
    return !isEmpty(this.emails) && !isEmpty(this.resourceSelectTemplate.value);
  }

  emptyFormField(): void {
    this.resourceSelectTemplate.select(-1);
  }

  show(): void {
    this.addResourceToMembersModalTemplate?.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog
        heading="Assign resource to members"
        @closed=${this.handleClosed}
      >
        <div class="emails">
          ${this.emails?.map((x: string) => {
            return html`<div>${x}</div>`;
          })}
        </div>
        <mwc-select label="Resources">
          ${this.resources?.map((x: ResourceResponse) => {
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
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
