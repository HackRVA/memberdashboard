// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  property,
} from "lit-element";

// membership
import { MemberService } from "../../service";
import "./member-list";
import { MemberLevel, MemberResponse } from "./types";

@customElement("members-page")
export class MembersPage extends LitElement {
  @property({ type: Array })
  members: MemberResponse[];

  @property({ type: Number })
  memberCount: number;

  memberService: MemberService = new MemberService();

  constructor() {
    super();
    this.members = [];
    this.memberCount = 0;
  }

  firstUpdated(): void {
    this.getMembers();
  }

  getMembers(): void {
    this.memberService.getMembers().subscribe({
      next: (result: any) => {
        this.members = result as MemberResponse[];
        this.memberCount = this.getActiveMembers().length;
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
