import { LitElement, html, TemplateResult } from "lit-element";

class MemberCount extends LitElement {
  memberCount: Number | null = null;

  firstUpdated(): void {
    this.memberCount = 2;
    this.requestUpdate();
  }
  render(): TemplateResult {
    return html` <h1>Member Count: ${this.memberCount}</h1> `;
  }
}

customElements.define("member-stats", MemberCount);
