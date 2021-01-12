import { LitElement, html, TemplateResult, customElement } from "lit-element";
import "./card-element";

@customElement('member-stats')
export class MemberCount extends LitElement {
  memberCount: Number | null = null;

  firstUpdated(): void {
    this.memberCount = 2;
    this.requestUpdate();
  }
  render(): TemplateResult {
    return html`
      <card-element><h1>Member Count: ${this.memberCount}</h1></card-element>
    `;
  }
}
