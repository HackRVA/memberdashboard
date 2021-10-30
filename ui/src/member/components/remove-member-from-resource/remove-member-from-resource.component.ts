// lit element
import { html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Dialog } from '@material/mwc-dialog';
import { Select } from '@material/mwc-select';
import { TextField } from '@material/mwc-textfield';

// memberdashboard
import { ResourceService } from '../../../resource/services/resource.service';
import { RemoveMemberResourceRequest } from '../../../resource/types/api/remove-member-resource-request';
import { isEmpty, showComponent } from '../../../shared/functions';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import { MemberResource } from '../../types/api/member-response';
import { Inject } from '../../../shared/di';

@customElement('remove-member-from-resource')
export class RemoveMemberFromResourceModal extends LitElement {
  @property({ type: String })
  email: string = '';

  @property({ type: Array })
  memberResources: Array<MemberResource> = [];

  toastMsg: ToastMessage;

  @Inject('resource')
  private resourceService: ResourceService;

  removeResourceFromMemberModalTemplate: Dialog;
  emailFieldTemplate: TextField;
  memberResourceSelectTemplate: Select;

  firstUpdated(): void {
    this.removeResourceFromMemberModalTemplate =
      this.shadowRoot.querySelector('mwc-dialog');
    this.emailFieldTemplate = this.shadowRoot.querySelector('mwc-textfield');
    this.memberResourceSelectTemplate =
      this.shadowRoot.querySelector('mwc-select');
  }

  show(): void {
    this.removeResourceFromMemberModalTemplate.show();
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.tryToRemoveMemberFromResource();
      this.emptyFormField();
      this.removeResourceFromMemberModalTemplate.close();
    } else {
      this.displayToastMsg('Hrmmmm');
    }
  }

  tryToRemoveMemberFromResource(): void {
    const request: RemoveMemberResourceRequest = {
      email: this.emailFieldTemplate.value,
      resourceID: this.memberResourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.removeMemberFromResource(request);
  }

  removeMemberFromResource(request: RemoveMemberResourceRequest): void {
    this.resourceService.removeMemberFromResource(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
      },
    });
  }

  fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
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

  emptyFormField(): void {
    this.memberResourceSelectTemplate.select(-1);
  }

  isValid(): boolean {
    return (
      !isEmpty(this.emailFieldTemplate.value) &&
      !isEmpty(this.memberResourceSelectTemplate.value)
    );
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Remove Resource" @closed=${this.handleClosed}>
        <mwc-textfield
          label="email"
          helper="Can't edit email"
          value=${this.email}
          readonly
        ></mwc-textfield>
        <mwc-select label="Resources">
          ${this.memberResources?.map((x: MemberResource) => {
            return html`
              <mwc-list-item value=${x.resourceID}> ${x.name} </mwc-list-item>
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
