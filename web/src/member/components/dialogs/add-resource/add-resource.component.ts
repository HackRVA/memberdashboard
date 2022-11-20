// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Select } from '@material/mwc-select';

// memberdashboard
import { ResourceService } from '../../../../resource/services/resource.service';
import { BulkAddMembersToResourceRequest } from '../../../../resource/types/api/bulk-add-members-to-resource-request';
import { ResourceResponse } from '../../../../resource/types/api/resource-response';
import { isEmpty } from '../../../../shared/functions';
import { addMembersToResourceStyle } from './add-resource.style';
import { Inject } from '../../../../shared/di/inject';
import { displayToast } from '../../../../shared/components/abstract-toast';
import { coreStyle } from '../../../../shared/styles';

@customElement('add-resource-form')
export class AddResourceForm extends LitElement {
  @property({ type: Array })
  emails: string[] = [];

  @property({ type: Function })
  closeHandler: () => void;

  resourceSelectTemplate: Select;

  @Inject('resource')
  private resourceService: ResourceService;

  resources: ResourceResponse[] = [];

  static get styles(): CSSResult[] {
    return [coreStyle, addMembersToResourceStyle];
  }

  firstUpdated(): void {
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
        // this.memberService.getMembers();
        displayToast('success');
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
      this.closeHandler();
      displayToast('success');
    } else {
      displayToast('Hrmmmmm');
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

  render(): TemplateResult {
    return html`
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
      <mwc-button slot="secondaryAction" @click=${this.closeHandler}>
        Cancel
      </mwc-button>
    `;
  }
}
