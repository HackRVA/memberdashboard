// lit element
import { LitElement, html, customElement, TemplateResult } from "lit-element";

// membership
import "../shared/card-element";
import "./resource-manager";

@customElement("resources-page")
export class ResourcesPage extends LitElement {
  render(): TemplateResult {
    return html`
      <card-element>
        <resource-manager> </resource-manager>
      </card-element>
    `;
  }
}
