import { LitElement, html, customElement, TemplateResult } from "lit-element";
import "./card-element";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-textfield";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-list";
import "@material/mwc-checkbox";
import { ResourceService } from "../service/resource.service";
import { isEmpty } from "../function";

const NOT_A_RESOURCE_ID = 0;

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: Array<ResourceService.ResourceResponse> | null = null;
  newAddress: string = "";
  newName: string = "";
  newID: number = 0;

  firstUpdated(): void {
    this.handleGetResources();
  }

  handleAddressChange(e: Event): void {
    this.newAddress = (e.target as EventTarget & { value: string }).value;
  }

  handleNameChange(e: Event): void {
    this.newName = (e.target as EventTarget & { value: string }).value;
  }

  handleOpenRegisterResource(): void {
    (this.shadowRoot?.querySelector("#register") as HTMLElement & {
      show: Function;
    }).show();
  }

  handleSubmitResource(isCreate: boolean): void {
    if (isCreate) {
      const request: ResourceService.RegisterResourceRequest = {
        name: this.newName,
        address: this.newAddress,
      };
      this.newID = NOT_A_RESOURCE_ID;
      this.newName = "";
      this.newAddress = "";
      this.handleRegisterResource(request);
    } else {
      const request: ResourceService.UpdateResourceRequest = {
        id: this.newID,
        name: this.newName,
        address: this.newAddress,
      };
      this.newID = NOT_A_RESOURCE_ID;
      this.newName = "";
      this.newAddress = "";
      this.handleUpdateResource(request);
    }
  }

  handleRegisterResource(
    request: ResourceService.RegisterResourceRequest
  ): void {
    this.resourceService.register(request).subscribe({
      complete: () => {
        this.handleGetResources();
        this.requestUpdate();
      },
    });
  }

  handleUpdateResource(request: ResourceService.UpdateResourceRequest): void {
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
        if ((result as { error: boolean; message: any }).error) {
          // this.onLoginComplete("Some error logging in");
          console.error("some error getting resources");
        } else {
          this.resources = result as ResourceService.ResourceResponse[];
          this.requestUpdate();
        }
      },
    });
  }

  handleDelete(resource: ResourceService.ResourceResponse): void {
    const request: ResourceService.RemoveResourceRequest = {
      id: resource.id,
    };
    this.resourceService.deleteResource(request).subscribe({
      complete: () => {
        this.handleGetResources();
        this.requestUpdate();
      },
    });
  }

  handleEdit(resource: ResourceService.UpdateResourceRequest): void {
    this.newAddress = resource.address;
    this.newName = resource.name;
    this.newID = resource.id || 0;
    this.requestUpdate();
    this.handleOpenRegisterResource();
  }

  emptyFormValues(): void {
    this.newID = NOT_A_RESOURCE_ID;
    this.newName = "";
    this.newAddress = "";
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
        @click=${this.emptyFormValues}
      >
        Cancel
      </mwc-button>
    </mwc-dialog>`;
  }

  resourceList(): TemplateResult | void {
    if (!this.resources) return;
    return html` <mwc-list>
      ${this.resources.map((x: ResourceService.ResourceResponse) => {
        return html`<mwc-list-item>
          ${x.name} ${x.address}
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
    return html` <card-element>
      <h1>Resource Manager</h1>
      <mwc-button
        @click=${this.handleOpenRegisterResource}
        dense
        unelevated
        label="create"
      ></mwc-button>

      <div>${this.resourceList()}</div>

      ${this.updateResourceDialog()}
    </card-element>`;
  }
}
