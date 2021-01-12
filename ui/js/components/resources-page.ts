import {
    LitElement,
    html,
    css,
    customElement,
    TemplateResult,
    CSSResult,
  } from "lit-element";
  
  @customElement("resources-page")
  export class ResourcesPage extends LitElement {
    static get styles(): CSSResult {
      return css``;
    }
    render(): TemplateResult {
      return html`
            <h1> Resources </h1>
      `;
    }
  }
  