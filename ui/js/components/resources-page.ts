import {
    LitElement,
    html,
    css,
    customElement,
    TemplateResult,
    CSSResult,
  } from "lit-element";
import "./resource-manager";

  
  @customElement("resources-page")
  export class ResourcesPage extends LitElement {
    static get styles(): CSSResult {
      return css``;
    }
    render(): TemplateResult {
      return html`
        <resource-manager> </resource-manager>
        `;
    }
  }
  