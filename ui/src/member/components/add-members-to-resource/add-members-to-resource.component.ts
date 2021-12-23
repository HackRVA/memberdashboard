// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Dialog } from '@material/mwc-dialog';
import { Select } from '@material/mwc-select';

// memberdashboard
import { ResourceService } from '../../../resource/services/resource.service';
import { BulkAddMembersToResourceRequest } from '../../../resource/types/api/bulk-add-members-to-resource-request';
import { ResourceResponse } from '../../../resource/types/api/resource-response';
import { isEmpty, showComponent } from '../../../shared/functions';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { addMembersToResourceStyle } from './add-members-to-resource.style';
import { Inject } from '../../../shared/di/inject';
import { IPopup } from './../../../shared/types/custom/ipop-up';

@customElement('add-members-to-resource')
export class AddMembersToResourceModal extends LitElement implements IPopup {
  @property({ type: Array })
  emails: string[] = [];

  addResourceToMembersModalTemplate: Dialog;
  resourceSelectTemplate: Select;

  @Inject('resource')
  private resourceService: ResourceService;

  resources: ResourceResponse[] = [];
  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [addMembersToResourceStyle];
  }

  firstUpdated(): void {
    this.addResourceToMembersModalTemplate =
      this.shadowRoot.querySelector('mwc-dialog');
    this.resourceSelectTemplate = this.shadowRoot.querySelector('mwc-select');

    this.getResources();
  }

  private getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result: ResourceResponse[]) => {
        this.resources = result;
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to get resources');
      },
    });
  }

  private displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  private tryToAddMembersToResource(): void {
    const request: BulkAddMembersToResourceRequest = {
      emails: this.emails,
      resourceID: this.resourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.bulkAddMembersToResource(request);
  }

  private bulkAddMembersToResource(
    request: BulkAddMembersToResourceRequest
  ): void {
    this.resourceService.bulkAddMembersToResource(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
      },
    });
  }

  private fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
  }

  private handleSubmit(): void {
    if (this.isValid()) {
      this.tryToAddMembersToResource();
      this.emptyFormField();
      this.addResourceToMembersModalTemplate.close();
    } else {
      this.displayToastMsg('Hrmmmmm');
    }
  }

  private handleClosed(event: CustomEvent): void {
    // temp hack to stop mwc-select from bubbling to mwc-dialog
    const tagName: string = (event.target as EventTarget & { tagName: string })
      .tagName;
    if (tagName === 'MWC-SELECT') {
      return;
    } else {
      this.emptyFormField();
    }
  }

  private isValid(): boolean {
    return !isEmpty(this.emails) && !isEmpty(this.resourceSelectTemplate.value);
  }

  private emptyFormField(): void {
    this.resourceSelectTemplate.select(-1);
  }

  public show(): void {
    this.addResourceToMembersModalTemplate?.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog
        heading="Assign resource to members"
        @closed=${this.handleClosed}
      >
        <div class="emails">
          ${this.emails?.map((x: string) => {
            return html`<div>${x}</div>`;
          })}
        </div>
        <mwc-select label="Resources">
          ${this.resources?.map((x: ResourceResponse) => {
            return html`
              <mwc-list-item value=${x.id}> ${x.name} </mwc-list-item>
            `;
          })}
        </mwc-select>
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
