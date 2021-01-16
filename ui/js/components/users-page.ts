import {
    LitElement,
    html,
    css,
    customElement,
    TemplateResult,
    CSSResult,
  } from "lit-element";
  
  @customElement("users-page")
  export class UsersPage extends LitElement {
    static get styles(): CSSResult {
      return css``;
    }
    render(): TemplateResult {
      return html`
            <h1> Users </h1>
      `;
    }
  }
  