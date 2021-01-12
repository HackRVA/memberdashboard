import {
    LitElement,
    html,
    css,
    customElement,
    TemplateResult,
    CSSResult,
  } from "lit-element";
import './card-element';
import './register-form'
  
  @customElement("home-page")
  export class HomePage extends LitElement {
    static get styles(): CSSResult {
      return css`
        login-container {
            display: grid;
            justify-content: center;
        }
      `;
    }
    render(): TemplateResult {
      return html`
            <card-element>
                <login-container> 
                    <register-form />
                </login-container>
            </card-element>
      `;
    }
  }