import { MemberService } from "./../../service/member.service";
import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";
import "../member-list";

@customElement("members-page")
export class MembersPage extends LitElement {
  @property({ type: Array })
  members: MemberService.MemberResponse[];

  @property({ type: Number })
  memberCount: number;

  memberService: MemberService = new MemberService();

  constructor() {
    super();
    this.members = [];
    this.memberCount = 0;
  }

  firstUpdated(): void {
    this.memberService.getMembers().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        this.members = result as MemberService.MemberResponse[];
        this.memberCount = this.members.length;
      },
    });
  }
  static get styles(): CSSResult {
    return css``;
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
