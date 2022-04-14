// lit element
import { html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Select } from '@material/mwc-select';
import { TextField } from '@material/mwc-textfield';

// memberdashboard
import { ResourceService } from '../../../../resource/services/resource.service';
import { RemoveMemberResourceRequest } from '../../../../resource/types/api/remove-member-resource-request';
import { isEmpty } from '../../../../shared/functions';
import { MemberResource } from '../../../types/api/member-response';
import { Inject } from '../../../../shared/di';
import { displayToast } from '../../../../shared/components/abstract-toast';
import { MemberManagerService } from '../../../services/member.service';

@customElement('remove-resource-form')
export class RemoveResourceForm extends LitElement {
  @property({ type: String })
  email: string = '';

  @property({ type: Array })
  memberResources: Array<MemberResource> = [];

  @property({ type: Function })
  closeHandler: () => void;

  @Inject('resource')
  private resourceService: ResourceService;

  @Inject('member-manager')
  private memberManagerService: MemberManagerService;

  emailFieldTemplate: TextField;
  memberResourceSelectTemplate: Select;

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot.querySelector('mwc-textfield');
    this.memberResourceSelectTemplate =
      this.shadowRoot.querySelector('mwc-select');
  }

  private handleSubmit(): void {
    if (this.isValid()) {
      this.tryToRemoveMemberFromResource();
      this.emptyFormField();
    } else {
      displayToast('Hrmmmm');
    }
  }

  private tryToRemoveMemberFromResource(): void {
    const request: RemoveMemberResourceRequest = {
      email: this.emailFieldTemplate.value,
      resourceID: this.memberResourceSelectTemplate.value,
    };
    this.emptyFormField();
    this.removeMemberFromResource(request);
  }

  private removeMemberFromResource(request: RemoveMemberResourceRequest): void {
    this.resourceService.removeMemberFromResource(request).subscribe({
      complete: () => {
        this.fireUpdatedEvent();
        this.memberManagerService.getMembers();
        displayToast('success');
      },
    });
  }

  private fireUpdatedEvent(): void {
    const updatedEvent = new CustomEvent('updated');
    this.dispatchEvent(updatedEvent);
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

  private emptyFormField(): void {
    this.memberResourceSelectTemplate.select(-1);
  }

  private isValid(): boolean {
    return (
      !isEmpty(this.emailFieldTemplate.value) &&
      !isEmpty(this.memberResourceSelectTemplate.value)
    );
  }

  render(): TemplateResult {
    return html`
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
      <mwc-button
        slot="primaryAction"
        @click=${() => {
          this.handleSubmit();
          this.closeHandler();
        }}
      >
        Submit
      </mwc-button>
      <mwc-button slot="secondaryAction" @click=${this.closeHandler}>
        Cancel
      </mwc-button>
    `;
  }
}
