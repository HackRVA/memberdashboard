// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
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
import { showComponent } from "../../function";
import { ResourceResponse, ResourceModalData } from "./types";
import { resourceManagerStyles } from "./styles/resource-manager.styles";
import "./modals";
import { ToastMessage } from "../shared/types";

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: Array<ResourceResponse> = [];

  resourceModalData: ResourceModalData;

  resourceName: string;
  resourceId: string;

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [resourceManagerStyles];
  }

  firstUpdated(): void {
    this.getResources();
  }

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result: any) => {
        this.resources = result as ResourceResponse[];
        this.requestUpdate();
      },
      error: () => {
        console.error("unable to get resources");
      },
    });
  }

  updateACLs(): void {
    this.resourceService.updateACLs().subscribe(() => {
      this.displayToastMsg("Successfully update ACL for all resource");
    });
  }

  removeACLs(): void {
    this.resourceService.removeACLs().subscribe(() => {
      this.displayToastMsg("Successfully remove ACL for all resource");
    });
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  openResourceWarningModal(resource: ResourceResponse): void {
    this.resourceName = resource.name;
    this.resourceId = resource.id;
    this.requestUpdate();
    showComponent("#resource-warning-modal", this.shadowRoot);
  }

  openRegisterResourceModal(): void {
    this.resourceModalData = Object.assign(
      {},
      {
        id: null,
        isEdit: false,
        resourceAddress: null,
        resourceName: null,
        isDefault: false,
      }
    );

    this.requestUpdate();
    showComponent("#resource-modal", this.shadowRoot);
  }

  openEditResourceModal(resource: ResourceResponse): void {
    this.resourceModalData = Object.assign(
      {},
      {
        id: resource.id,
        isEdit: true,
        resourceName: resource.name,
        resourceAddress: resource.address,
        isDefault: resource.isDefault,
      }
    );
    this.requestUpdate();
    showComponent("#resource-modal", this.shadowRoot);
  }

  refreshResources(): void {
    this.getResources();
    this.requestUpdate();
  }

  displayResources(): TemplateResult {
    return html`
      ${this.resources?.map((x: ResourceResponse) => {
        return html`
          <tr>
            <td>${x.name} ${x.isDefault ? "(default)" : ""}</td>
            <td>${x.address}</td>
            <td>
              <mwc-button
                @click="${() => this.openEditResourceModal(x)}"
                label="Edit"
              >
              </mwc-button>
              <mwc-button
                class="remove"
                @click="${() => this.openResourceWarningModal(x)}"
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
              class="mr-8"
              @click=${this.updateACLs}
              label="Update ACLs"
              dense
              unelevated
            ></mwc-button>
            <mwc-button
              class="mr-32 remove"
              @click=${this.removeACLs}
              label="Delete ACLs"
              dense
              unelevated
            ></mwc-button>
            <mwc-button
              class="create-resource"
              @click=${this.openRegisterResourceModal}
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
        @updated=${this.refreshResources}
        .resourceModalData=${this.resourceModalData}
      >
      </resource-modal>
      <warning-modal
        id="resource-warning-modal"
        .resourceName=${this.resourceName}
        .resourceId=${this.resourceId}
        @updated=${this.refreshResources}
      >
      </warning-modal>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
