// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import { MemberService } from "../../service";
import "./member-list";
import { MemberLevel, MemberResponse } from "./types";
import "../shared/card-element";
import { membersPageStyles } from "./styles/members-page-styles";

@customElement("members-page")
export class MembersPage extends LitElement {
  members: MemberResponse[];
  memberCount: number;
  showSpinner: boolean = true;

  memberService: MemberService = new MemberService();

  static get styles(): CSSResult[] {
    return [membersPageStyles];
  }

  firstUpdated(): void {
    this.getMembers();
  }

  getMembers(): void {
    this.memberService.getMembers().subscribe({
      next: (result: MemberResponse[]) => {
        this.showSpinner = false;
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

  displayMemberList(): TemplateResult {
    if (this.showSpinner) {
      return html`<mwc-circular-progress
        indeterminate
      ></mwc-circular-progress>`;
    }

    return html`
      <member-list
        .members=${this.members}
        .memberCount=${this.memberCount}
      ></member-list>
    `;
  }

  render(): TemplateResult {
    return html`
      <card-element class="text-center">
        ${this.displayMemberList()}
      </card-element>
    `;
  }
}
