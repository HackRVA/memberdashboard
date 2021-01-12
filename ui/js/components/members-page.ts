import {
    LitElement,
    html,
    css,
    customElement,
    TemplateResult,
    CSSResult,
  } from "lit-element";
  
  @customElement("members-page")
  export class MembersPage extends LitElement {
    static get styles(): CSSResult {
      return css``;
    }
    render(): TemplateResult {
      return html`
            <h1> Members </h1>
      `;
    }
  }
  