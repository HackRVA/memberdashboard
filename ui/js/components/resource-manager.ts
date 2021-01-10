import { LitElement, html, customElement, TemplateResult } from "lit-element";
import "./card-element";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-textfield";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-list";
import "@material/mwc-checkbox";
import { ResourceService } from "../service/Resource";

const NOT_A_RESOURCE_ID = 0;

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: Array<{ id: number; address: string; name: string }> | null = null;
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

  handleRegisterResource(): void {
    this.resourceService
      .register({
        id: this.newID,
        name: this.newName,
        address: this.newAddress,
      })
      .subscribe();

    this.newID = NOT_A_RESOURCE_ID;
    this.newName = "";
    this.newAddress = "";
  }

  handleGetResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result) => {
        if ((result as { error: boolean; message: any }).error) {
          // this.onLoginComplete("Some error logging in");
          console.error("some error getting resources");
        } else {
          this.resources = result as any;
          this.requestUpdate();
        }
      },
      // complete: () => this.onLoginComplete("Success!"),
    });
  }

  handleDelete(resource: ResourceService.ResourceRequest): void {
    this.resourceService.deleteResource(resource).subscribe({
      complete: () => {
        this.handleGetResources();
        this.requestUpdate();
      },
    });
  }

  handleEdit(resource: ResourceService.ResourceRequest): void {
    this.newAddress = resource.address;
    this.newName = resource.name;
    this.newID = resource.id || 0;
    this.requestUpdate();
    this.handleOpenRegisterResource();
  }

  updateResourceDialog(): TemplateResult {
    return html`<mwc-dialog id="register">
      <div>Update/Register a Resource?</div>

      <mwc-textfield
        @change=${this.handleNameChange}
        id="newResourceName"
        label="name"
        helper="name of device"
      ></mwc-textfield>
      <mwc-textfield
        @change=${this.handleAddressChange}
        id="newResourceAddress"
        label="address"
        helper="address on the network"
      ></mwc-textfield>

      <mwc-button
        @click=${this.handleRegisterResource}
        slot="primaryAction"
        dialogAction="discard"
      >
        Submit
      </mwc-button>
      <mwc-button slot="secondaryAction" dialogAction="cancel">
        Cancel
      </mwc-button>
    </mwc-dialog>`;
  }

  resourceList(): TemplateResult | void {
    if (!this.resources) return;
    return html` <mwc-list>
      ${this.resources.map((x: ResourceService.ResourceRequest) => {
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
