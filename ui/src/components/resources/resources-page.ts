// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import "../shared/card-element";
import "./resource-manager";
import { ResourceResponse } from "./types/resource.interface";
import { ResourceService } from "./../../service/resource.service";
import "../shared/loading-content";
import { coreStyles } from "./../shared/styles/core-styles";
import { resourcesPageStyles } from "./styles/resources-page-styles";

@customElement("resources-page")
export class ResourcesPage extends LitElement {
  resourceService: ResourceService = new ResourceService();
  resources: ResourceResponse[];
  finishedLoading: boolean = false;

  static get styles(): CSSResult[] {
    return [resourcesPageStyles, coreStyles];
  }

  firstUpdated(): void {
    this.getResources();
  }

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result: ResourceResponse[]) => {
        this.finishedLoading = true;
        this.resources = result;
        this.requestUpdate();
      },
      error: () => {
        console.error("unable to get resources");
      },
    });
  }

  render(): TemplateResult {
    return html`
      <card-element class="center-text">
        <loading-content .finishedLoading=${this.finishedLoading}>
          <resource-manager .resources=${this.resources}> </resource-manager>
        </loading-content>
      </card-element>
    `;
  }
}
