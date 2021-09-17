// lit element
import { html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Dialog } from '@material/mwc-dialog';

// memberdashboard
import { showComponent } from '../../../shared/functions';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { ResourceService } from '../../services/resource.service';
import { RemoveResourceRequest } from '../../types/api/remove-resource-request';

@customElement('resource-warning')
export class ResourceWarning extends LitElement {
  @property({ type: String })
  resourceName: string;

  @property({ type: String })
  resourceId: string;

  toastMsg: ToastMessage;

  resourceWarningModalTemplate: Dialog;

  resourceService: ResourceService = new ResourceService();

  firstUpdated(): void {
    this.resourceWarningModalTemplate =
      this.shadowRoot.querySelector('mwc-dialog');
  }

  show(): void {
    this.resourceWarningModalTemplate.show();
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  private tryToRemoveResource(): void {
    const request: RemoveResourceRequest = {
      id: this.resourceId,
    };

    this.handleRemoveResource(request);
  }

  private handleRemoveResource(request: RemoveResourceRequest): void {
    this.resourceService.deleteResource(request).subscribe({
      complete: () => {
        this.displayToastMsg(
          'Success, but I hope you know what you are doing :3'
        );
        this.fireUpdatedEvent();
        this.resourceWarningModalTemplate.close();
      },
    });
  }

  private fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
  }

  private handleSubmit(): void {
    this.tryToRemoveResource();
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Warning">
        <div>
          Are you sure you want to remove
          <strong>${this.resourceName}</strong> ? <br />
          Members will no longer have access to this resource if you remove it.
        </div>
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
