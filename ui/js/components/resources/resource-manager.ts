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
import { showComponent } from "../../function";
import { ResourceResponse, ResourceModalData } from "./types";
import "./modals";

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: Array<ResourceResponse> = [];

  resourceModalData: ResourceModalData;

  resourceName: string;
  resourceId: string;

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
    this.getResources();
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
      ${this.resources.map((x: ResourceResponse) => {
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
    `;
  }
}
