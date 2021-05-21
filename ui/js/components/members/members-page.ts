// lit element
import { LitElement, html, customElement, TemplateResult } from "lit-element";

// membership
import { MemberService } from "../../service";
import "./member-list";
import { MemberLevel, MemberResponse } from "./types";

@customElement("members-page")
export class MembersPage extends LitElement {
  members: MemberResponse[];
  memberCount: number;

  memberService: MemberService = new MemberService();

  firstUpdated(): void {
    this.getMembers();
  }

  getMembers(): void {
    this.memberService.getMembers().subscribe({
      next: (result: MemberResponse[]) => {
        this.members = result;
        this.memberCount = this.getActiveMembers().length;
        this.requestUpdate();
      },
      error: () => {
        console.error("unable to get members");
      },
    });
  }

  getActiveMembers(): MemberResponse[] {
    return this.members.filter(
      (x: MemberResponse) => x.memberLevel !== MemberLevel.inactive
    );
  }

  render(): TemplateResult {
    return html`
      <member-list
        .members=${this.members}
        .memberCount=${this.memberCount}
      ></member-list>
    `;
  }
}
