// lit element
import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";

// membership
import { MemberService } from "../../service/member.service";
import "./member-list";
import { MemberResponse } from "./types";

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
    this.memberService.getMembers().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        this.members = result as MemberResponse[];
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
