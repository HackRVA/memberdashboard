// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";

// memberdashboard
import "../shared/card-element";
import { ResourceService } from "../../service";
import { showComponent } from "../../function";
import { ResourceResponse, ResourceModalData } from "./types";
import { coreStyles } from "./../shared/styles/core-styles";
import { resourceManagerStyles } from "./styles/resource-manager-styles";
import "./modals";
import { ToastMessage } from "../shared/types";

@customElement("resource-manager")
export class ResourceManager extends LitElement {
  @property({ type: Array })
  resources: ResourceResponse[] = [];

  @property({ type: Number })
  resourceCount: number = 0;

  resourceService: ResourceService = new ResourceService();

  resourceModalData: ResourceModalData;

  resourceName: string;
  resourceId: string;

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [resourceManagerStyles, coreStyles];
  }

  private getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (response: ResourceResponse[]) => {
        this.resources = response;
        this.resourceCount = response.length;
        this.requestUpdate();
      },
      error: () => {
        console.error("unable to get resources");
      },
    });
  }

  private updateACLs(): void {
    this.resourceService.updateACLs().subscribe(() => {
      this.displayToastMsg("Successfully update ACL for all resource");
    });
  }

  private removeACLs(): void {
    this.resourceService.removeACLs().subscribe(() => {
      this.displayToastMsg("Successfully remove ACL for all resource");
    });
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  private openResourceWarningModal(resource: ResourceResponse): void {
    this.resourceName = resource.name;
    this.resourceId = resource.id;
    this.requestUpdate();
    showComponent("#resource-warning-modal", this.shadowRoot);
  }

  private openRegisterResourceModal(): void {
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

  private openEditResourceModal(resource: ResourceResponse): void {
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

  private refreshResources(): void {
    this.getResources();
    this.requestUpdate();
  }

  private displayOnlineState(lastHeartBeat: string): TemplateResult {
    const lastActiveTime: number = new Date(lastHeartBeat).getTime();
    const currentTime: number = new Date().getTime();

    if (this.isTimeWitihinActiveRange(lastActiveTime, currentTime)) {
      return html` <span class="online"> Online </span> `;
    } else {
      return html`<span class="offline"> Offline </span>`;
    }
  }

  private isTimeWitihinActiveRange(
    lastActiveTime: number,
    currentTime: number
  ): boolean {
    const thirtyMinsInMS: number = 30 * 60 * 1000;
    const remainingTime: number = currentTime - lastActiveTime;
    if (remainingTime > 0 && remainingTime <= thirtyMinsInMS) {
      return true;
    }

    return false;
  }

  private displayResources(): TemplateResult {
    return html`
      ${this.resources?.map((x: ResourceResponse) => {
        return html`
          <tr>
            <td>${x.name} ${x.isDefault ? "(default)" : ""}</td>
            <td>${x.address}</td>
            <td>${this.displayOnlineState(x.lastHeartBeat)}</td>
            <td>
              <mwc-button
                @click="${() => this.openEditResourceModal(x)}"
                label="Edit"
              >
              </mwc-button>
              <mwc-button
                class="destructive-button"
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
          <h3>Number of resources: ${this.resourceCount}</h3>
          <div>
            <span class="update-acls">
              <mwc-button
                @click=${this.updateACLs}
                label="Update ACLs"
                dense
                unelevated
              ></mwc-button>
              <div class="note">
                Update ACLs adds members and does not remove members
              </div>
            </span>
            <mwc-button
              class="remove-acls"
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
            <th>Status</th>
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
