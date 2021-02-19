// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
  css,
} from "lit-element";

// material
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-textfield";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-list";
import "@material/mwc-checkbox";

// membership
import "../shared/card-element";
import { ResourceService } from "../../service";
import { isEmpty, showComponent } from "../../function";
import {
  ResourceResponse,
  RegisterResourceRequest,
  UpdateResourceRequest,
  RemoveResourceRequest,
  ResourceModalData,
} from "./types";
import "./modals";

const NOT_A_RESOURCE_ID = "";

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: Array<ResourceResponse> = [];
  newAddress: string = "";
  newName: string = "";
  newID: string = "";
  newIsDefault: boolean = false;

  resourceModalData: ResourceModalData;

  static get styles(): CSSResult {
    return css`
      .resource-container {
        display: grid;
        align-items: center;
        margin: 44px;
      }

      .resource-header {
        display: inherit;
        grid-template-columns: 1fr 1fr;
        align-items: center;
      }

      .button-container {
        justify-self: end;
      }

      td,
      th {
        text-align: left;
        padding: 8px;
        font-size: 20px;
        border: 1px solid #e1e1e1;
        max-width: 320px;
      }
      table {
        margin-top: 24px;
        border-spacing: 0px;
      }

      .remove {
        --mdc-theme-primary: #e9437a;
      }
    `;
  }

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
    this.newIsDefault = (e.target as EventTarget & {
      checked: boolean;
    }).checked;
  }

  openRegisterResource(): void {
    this.resourceModalData = Object.assign(
      {},
      {
        isEdit: false,
        resourceAddress: null,
        resourceName: null,
        isDefault: false,
      }
    );

    this.requestUpdate();
    showComponent("#resource-modal", this.shadowRoot);
  }

  handleSubmitResource(isCreate: boolean): void {
    if (isCreate) {
      const request: RegisterResourceRequest = {
        name: this.newName,
        address: this.newAddress,
        isDefault: this.newIsDefault,
      };
      this.emptyFormValues();
      this.handleRegisterResource(request);
    } else {
      const request: UpdateResourceRequest = {
        id: this.newID,
        name: this.newName,
        address: this.newAddress,
        isDefault: this.newIsDefault,
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
    this.newIsDefault = resource.isDefault;
    this.requestUpdate();
    this.openRegisterResource();
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

    return html`
      <mwc-dialog id="register" heading="Update/Register a Resource?">
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
            ?checked=${this.newIsDefault}
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
      </mwc-dialog>
    `;
  }

  displayResources(): TemplateResult {
    return html`
      ${this.resources.map((x: ResourceResponse) => {
        return html`
          <tr>
            <td>${x.name} ${x.isDefault ? "(default)" : ""}</td>
            <td>${x.address}</td>
            <td>
              <mwc-button @click="${() => this.handleEdit(x)}" label="Edit">
              </mwc-button>
              <mwc-button
                class="remove"
                @click="${() => this.handleDelete(x)}"
                label="Delete"
              >
              </mwc-button>
            </td>
          </tr>
        `;
      })}
    `;
  }

  render(): TemplateResult {
    return html`
      <div class="resource-container">
        <div class="resource-header">
          <h1>Resources</h1>
          <div class="button-container">
            <mwc-button
              class="create-resource"
              @click=${this.openRegisterResource}
              dense
              unelevated
              label="create"
            ></mwc-button>
          </div>
        </div>
        <table>
          <tr>
            <th>Name</th>
            <th>Address</th>
            <th>Actions</th>
          </tr>
          ${this.displayResources()}
        </table>
      </div>
      <resource-modal
        id="resource-modal"
        .resourceModalData=${this.resourceModalData}
      >
      </resource-modal>
    `;
  }
}
