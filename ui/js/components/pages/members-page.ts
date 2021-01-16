import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import "./../member-stats";

@customElement("members-page")
export class MembersPage extends LitElement {
  static get styles(): CSSResult {
    return css``;
  }
  render(): TemplateResult {
    return html` <member-stats></member-stats> `;
  }
}
