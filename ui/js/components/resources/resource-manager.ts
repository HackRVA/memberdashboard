// lit element
import { LitElement, html, customElement, TemplateResult } from "lit-element";

// material
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-textfield";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-list";
import "@material/mwc-checkbox";

// membership
import "../shared/card-element";
import { ResourceService } from "../../service/resource.service";
import { isEmpty, showComponent } from "../../function";
import {
  ResourceResponse,
  RegisterResourceRequest,
  UpdateResourceRequest,
  RemoveResourceRequest,
} from "./types";

const NOT_A_RESOURCE_ID = "";

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: Array<ResourceResponse> = [];
  newAddress: string = "";
  newName: string = "";
  newID: string = "";
  newIsDefault: boolean = false;

  firstUpdated(): void {
    this.handleGetResources();
  }

  handleAddressChange(e: Event): void {
    this.newAddress = (e.target as EventTarget & { value: string }).value;
  }

  handleNameChange(e: Event): void {
    this.newName = (e.target as EventTarget & { value: string }).value;
  }

  handleIsDefaultChange(e: Event): void {
    this.newIsDefault = (e.target as EventTarget & { checked: boolean }).checked;
  }

  handleOpenRegisterResource(): void {
    showComponent("#register", this.shadowRoot);
  }

  handleSubmitResource(isCreate: boolean): void {
    if (isCreate) {
      const request: RegisterResourceRequest = {
        name: this.newName,
        address: this.newAddress,
        is_default: this.newIsDefault
      };
      this.emptyFormValues();
      this.handleRegisterResource(request);
    } else {
      const request: UpdateResourceRequest = {
        id: this.newID,
        name: this.newName,
        address: this.newAddress,
        is_default: this.newIsDefault,
      };
      this.emptyFormValues();
      this.handleUpdateResource(request);
    }
  }

  handleRegisterResource(request: RegisterResourceRequest): void {
    this.resourceService.register(request).subscribe({
      complete: () => {
        this.handleGetResources();
        this.requestUpdate();
      },
    });
  }

  handleUpdateResource(request: UpdateResourceRequest): void {
    this.resourceService.updateResource(request).subscribe({
      complete: () => {
        this.handleGetResources();
        this.requestUpdate();
      },
    });
  }

  handleGetResources(): void {
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

  handleDelete(resource: ResourceResponse): void {
    const request: RemoveResourceRequest = {
      id: resource.id,
    };
    this.resourceService.deleteResource(request).subscribe({
      complete: () => {
        this.handleGetResources();
        this.requestUpdate();
      },
    });
  }

  handleEdit(resource: UpdateResourceRequest): void {
    this.newAddress = resource.address;
    this.newName = resource.name;
    this.newID = resource.id;
    this.newIsDefault = resource.is_default;
    this.requestUpdate();
    this.handleOpenRegisterResource();
  }

  emptyFormValues(): void {
    this.newID = NOT_A_RESOURCE_ID;
    this.newName = "";
    this.newAddress = "";
    this.newIsDefault = false;
  }

  emptyFormValuesOnClosed(): void {
    this.emptyFormValues();
    this.requestUpdate();
  }

  updateResourceDialog(): TemplateResult {
    const isCreate: boolean = isEmpty(this.newName) && isEmpty(this.newAddress);

    return html`<mwc-dialog id="register">
      <div>Update/Register a Resource?</div>

      <mwc-textfield
        @change=${this.handleNameChange}
        value=${this.newName}
        id="newResourceName"
        label="name"
        helper="Name of device"
      ></mwc-textfield>
      <mwc-textfield
        @change=${this.handleAddressChange}
        value=${this.newAddress}
        id="newResourceAddress"
        label="address"
        helper="Address on the network"
      ></mwc-textfield>
      <mwc-formfield label="Default">
        <mwc-checkbox
          @change=${this.handleIsDefaultChange}
          checked=${this.newIsDefault}
          value="true"
          id="newResourceIsDefault"
          helper="Default Resource for New Users"
        ></mwc-checkbox>
      </mwc-formfield>

      <mwc-button
        @click=${() => this.handleSubmitResource(isCreate)}
        slot="primaryAction"
        dialogAction="ok"
      >
        Submit
      </mwc-button>
      <mwc-button
        slot="secondaryAction"
        dialogAction="cancel"
        @click=${this.emptyFormValuesOnClosed}
      >
        Cancel
      </mwc-button>
    </mwc-dialog>`;
  }

  resourceList(): TemplateResult | void {
    if (!this.resources) return;
    return html` <mwc-list>
      ${this.resources.map((x: ResourceResponse) => {
        return html`<mwc-list-item>
          ${x.name} ${x.address} ${x.is_default ? '(assigned by default)' : ''}
          <mwc-button
            @click="${() => this.handleDelete(x)}"
            label="delete"
          ></mwc-button>
          <mwc-button
            @click="${() => this.handleEdit(x)}"
            label="edit"
          ></mwc-button>
        </mwc-list-item>`;
      })}
    </mwc-list>`;
  }

  render(): TemplateResult {
    return html`
      <div>
        <mwc-button
          @click=${this.handleOpenRegisterResource}
          dense
          unelevated
          label="create"
        ></mwc-button>

        <div>${this.resourceList()}</div>

        ${this.updateResourceDialog()}
      </div>
    `;
  }
}
