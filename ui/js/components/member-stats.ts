import { LitElement, html, TemplateResult } from "lit-element";
import "./card-element";

class MemberCount extends LitElement {
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

customElements.define("member-stats", MemberCount);
