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
import { coreStyles } from "./../shared/styles/core-styles";
import { membersPageStyles } from "./styles/members-page-styles";
import "../shared/loading-content";

@customElement("members-page")
export class MembersPage extends LitElement {
  members: MemberResponse[];
  memberCount: number;
  finishedLoading: boolean = false;

  memberService: MemberService = new MemberService();

  static get styles(): CSSResult[] {
    return [membersPageStyles, coreStyles];
  }

  firstUpdated(): void {
    this.getMembers();
  }

  getMembers(): void {
    this.memberService.getMembers().subscribe({
      next: (result: MemberResponse[]) => {
        this.finishedLoading = true;
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
      <card-element class="center-text">
        <loading-content .finishedLoading=${this.finishedLoading}>
          <member-list
            .members=${this.members}
            .memberCount=${this.memberCount}
          ></member-list>
        </loading-content>
      </card-element>
    `;
  }
}
