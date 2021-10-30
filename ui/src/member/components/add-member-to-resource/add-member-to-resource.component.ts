// lit element
import { customElement, property } from 'lit/decorators.js';
import { html, LitElement, TemplateResult } from 'lit';

// material
import { TextField } from '@material/mwc-textfield/mwc-textfield';
import { Dialog } from '@material/mwc-dialog';
import { Select } from '@material/mwc-select';

// memberdashboard
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { ResourceService } from '../../../resource/services/resource.service';
import { ResourceResponse } from '../../../resource/types/api/resource-response';
import { BulkAddMembersToResourceRequest } from '../../../resource/types/api/bulk-add-members-to-resource-request';
import { isEmpty, showComponent } from '../../../shared/functions';
import { Inject } from '../../../shared/di/inject';

@customElement('add-member-to-resource')
export class AddMemberToResourceModal extends LitElement {
  @property({ type: String })
  email: string = '';

  @Inject('resource')
  private resourceService: ResourceService;

  resources: ResourceResponse[] = [];
  toastMsg: ToastMessage;

  addResourceToMemberModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  resourceSelectTemplate: Select;

  firstUpdated(): void {
    this.addResourceToMemberModalTemplate =
      this.shadowRoot.querySelector('mwc-dialog');
    this.emailFieldTemplate = this.shadowRoot.querySelector('mwc-textfield');
    this.resourceSelectTemplate = this.shadowRoot.querySelector('mwc-select');
    this.getResources();
  }

  show(): void {
    this.addResourceToMemberModalTemplate?.show();
  }

  getResources(): void {
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

  handleClosed(event: CustomEvent): void {
    // temp hack to stop mwc-select from bubbling to mwc-dialog
    const tagName: string = (event.target as EventTarget & { tagName: string })
      .tagName;
    if (tagName === 'MWC-SELECT') {
      return;
    } else {
      this.emptyFormField();
    }
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.tryToAddMemberToResource();
      this.emptyFormField();
      this.addResourceToMemberModalTemplate.close();
    } else {
      this.displayToastMsg('Hrmmmmm');
    }
  }

  tryToAddMemberToResource(): void {
    const request: BulkAddMembersToResourceRequest = {
      emails: [this.emailFieldTemplate.value],
      resourceID: this.resourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.addMemberToResource(request);
  }

  addMemberToResource(request: BulkAddMembersToResourceRequest): void {
    this.resourceService.bulkAddMembersToResource(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
      },
    });
  }

  fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  emptyFormField(): void {
    this.resourceSelectTemplate.select(-1);
  }

  isValid(): boolean {
    return (
      !isEmpty(this.emailFieldTemplate.value) &&
      !isEmpty(this.resourceSelectTemplate.value)
    );
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Add Resource" @closed=${this.handleClosed}>
        <mwc-textfield
          label="email"
          helper="Can't edit email"
          value=${this.email}
          readonly
        ></mwc-textfield>
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
