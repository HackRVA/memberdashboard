// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// polymer
import '@polymer/paper-tooltip';

// memberdashboard
import '../resource-info';
import '../resource-warning';
import { showComponent } from '../../../shared/functions';
import { coreStyle } from '../../../shared/styles';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { ResourceService } from '../../services/resource.service';
import { ResourceResponse } from '../../types/api/resource-response';
import { ResourceModalData } from '../../types/custom/resource-modal-data';
import { resourceManagerStyle } from './resource-manager.style';

@customElement('resource-manager')
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
    return [resourceManagerStyle, coreStyle];
  }

  private getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (response: ResourceResponse[]) => {
        this.resources = response;
        this.resourceCount = response.length;
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to get resources');
      },
    });
  }

  private updateACLs(): void {
    this.resourceService.updateACLs().subscribe(() => {
      this.displayToastMsg('Successfully update ACL for all resource');
    });
  }

  private removeACLs(): void {
    this.resourceService.removeACLs().subscribe(() => {
      this.displayToastMsg('Successfully remove ACL for all resource');
    });
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  private openResourceWarningModal(resource: ResourceResponse): void {
    this.resourceName = resource.name;
    this.resourceId = resource.id;
    this.requestUpdate();
    showComponent('#resource-warning-modal', this.shadowRoot);
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
    showComponent('#resource-modal', this.shadowRoot);
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
    showComponent('#resource-modal', this.shadowRoot);
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
            <td>${x.name} ${x.isDefault ? '(default)' : ''}</td>
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
            <mwc-button
              id="update-acls"
              @click=${this.updateACLs}
              dense
              unelevated
              label="update-acls"
            ></mwc-button>
            <paper-tooltip for="update-acls" animation-delay="0">
              Update ACLs adds members and does not remove members
            </paper-tooltip>
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
      <resource-info
        id="resource-modal"
        @updated=${this.refreshResources}
        .resourceModalData=${this.resourceModalData}
      >
      </resource-info>
      <resource-warning
        id="resource-warning-modal"
        .resourceName=${this.resourceName}
        .resourceId=${this.resourceId}
        @updated=${this.refreshResources}
      >
      </resource-warning>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
