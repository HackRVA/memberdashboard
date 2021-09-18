// lit element
import { customElement } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../../../shared/components/md-card';
import '../resource-manager';
import '../../../shared/components/loading-content';
import { coreStyle } from '../../../shared/styles';
import { ResourceService } from '../../services/resource.service';
import { ResourceResponse } from '../../types/api/resource-response';

@customElement('resource-page')
export class ResourcePage extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: ResourceResponse[];
  resourceCount: number;
  finishedLoading: boolean = false;

  static get styles(): CSSResult[] {
    return [coreStyle];
  }

  firstUpdated(): void {
    this.getResources();
  }

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (response: ResourceResponse[]) => {
        this.finishedLoading = true;
        this.resources = response;
        this.resourceCount = response.length;
        this.requestUpdate();
      },
      error: () => {
        console.error('unable to get resources');
      },
    });
  }

  render(): TemplateResult {
    return html`
      <md-card>
        <loading-content .finishedLoading=${this.finishedLoading}>
          <resource-manager
            .resources=${this.resources}
            .resourceCount=${this.resourceCount}
          >
          </resource-manager>
        </loading-content>
      </md-card>
    `;
  }
}
