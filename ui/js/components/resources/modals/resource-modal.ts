// lit element
import {
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// material
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-textfield";
import "@material/mwc-checkbox";
import { TextField } from "@material/mwc-textfield/mwc-textfield";
import { Dialog } from "@material/mwc-dialog";
import { Checkbox } from "@material/mwc-checkbox";

// membership
import { ResourceService } from "./../../../service";
import {
  RegisterResourceRequest,
  UpdateResourceRequest,
  ResourceModalData,
} from "../types";
import "../../shared/toast-msg";
import { ToastMessage } from "../../shared/types";
import { showComponent } from "../../../function";

@customElement("resource-modal")
export class ResourceModal extends LitElement {
  @property({ type: Object })
  resourceModalData: ResourceModalData;

  toastMsg: ToastMessage;

  resourceModalTemplate: Dialog;
  resourceNameFieldTemplate: TextField;
  resourceAddressFieldTemplate: TextField;
  defaultResourceTemplate: Checkbox;

  resourceService: ResourceService = new ResourceService();

  firstUpdated(): void {
    this.resourceModalTemplate = this.shadowRoot.querySelector("mwc-dialog");
    this.resourceNameFieldTemplate = this.shadowRoot.querySelector(
      "#resource-name"
    );
    this.resourceAddressFieldTemplate = this.shadowRoot.querySelector(
      "#resource-address"
    );
    this.defaultResourceTemplate = this.shadowRoot.querySelector(
      "mwc-checkbox"
    );
  }

  updated(): void {
    if (this.resourceModalData?.isEdit) {
      this.resourceNameFieldTemplate.value = this.resourceModalData.resourceName;
      this.resourceAddressFieldTemplate.value = this.resourceModalData.resourceAddress;
      this.defaultResourceTemplate.checked = this.resourceModalData.isDefault;
    }
  }

  show(): void {
    this.resourceModalTemplate.show();
  }

  private trytoRegisterResource(): void {
    const request: RegisterResourceRequest = {
      name: this.resourceNameFieldTemplate.value,
      address: this.resourceAddressFieldTemplate.value,
      isDefault: this.defaultResourceTemplate.checked,
    };

    this.handleRegisterResource(request);
  }

  private trytoUpdateResource(): void {
    const request: UpdateResourceRequest = {
      id: this.resourceModalData.id,
      name: this.resourceNameFieldTemplate.value,
      address: this.resourceAddressFieldTemplate.value,
      isDefault: this.defaultResourceTemplate.checked,
    };

    this.handleUpdateResource(request);
  }

  private handleRegisterResource(request: RegisterResourceRequest): void {
    this.resourceService.register(request).subscribe({
      complete: () => {
        this.displayToastMsg("Success");
        this.emptyFormField();
        this.fireUpdatedEvent();
        this.resourceModalTemplate.close();
      },
    });
  }

  private handleUpdateResource(request: UpdateResourceRequest): void {
    this.resourceService.updateResource(request).subscribe({
      complete: () => {
        this.displayToastMsg("Success");
        this.emptyFormField();
        this.fireUpdatedEvent();
        this.resourceModalTemplate.close();
      },
    });
  }

  private fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent("updated");
    this.dispatchEvent(updatedEvent);
  }

  private handleSubmit(): void {
    if (this.isValid()) {
      if (this.resourceModalData.isEdit) {
        this.trytoUpdateResource();
      } else {
        this.trytoRegisterResource();
      }
    } else {
      this.displayToastMsg(
        "Hrmmm, are you sure everything in the form is correct?"
      );
    }
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  private emptyFormField(): void {
    this.resourceNameFieldTemplate.value = "";
    this.resourceAddressFieldTemplate.value = "";
    this.defaultResourceTemplate.checked = false;
  }

  private handleClosed(): void {
    this.emptyFormField();
  }

  private isValid(): boolean {
    return (
      this.resourceNameFieldTemplate.validity.valid &&
      this.resourceAddressFieldTemplate.validity.valid
    );
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog
        heading="Update/Register a Resource"
        @closed=${this.handleClosed}
      >
        <mwc-textfield
          required
          id="resource-name"
          label="name"
          helper="Name of device"
        ></mwc-textfield>
        <mwc-textfield
          required
          id="resource-address"
          label="address"
          helper="Address on the network"
        ></mwc-textfield>
        <mwc-formfield label="Default">
          <mwc-checkbox></mwc-checkbox>
        </mwc-formfield>
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
