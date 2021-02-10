// lit element
import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import "./resource-manager";

// membership
import "../shared/card-element";

@customElement("resources-page")
export class ResourcesPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html`
      <card-element>
        <h1>Resources</h1>
        <resource-manager> </resource-manager>
      </card-element>
    `;
  }
}
