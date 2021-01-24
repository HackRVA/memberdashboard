import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  css,
  CSSResult,
} from "lit-element";
import { MemberService } from "../service/member.service";
import "./card-element";

@customElement("member-list")
export class MemberList extends LitElement {
  memberCount: number = 0;
  members: MemberService.MemberResponse[] = [];
  memberService: MemberService = new MemberService();

  static get styles(): CSSResult {
    return css`
      h1 {
        margin-top: 0px;
      }
      .member-container {
        display: grid;
        justify-content: center;
        text-align: center;
        margin: 44px;
      }
      .name {
        text-transform: capitalize;
      }
      td,
      th {
        text-align: left;
        padding: 8px;
        font-size: 24px;
        border: 1px solid #e1e1e1;
        width: 320px;
      }
      table {
        margin-top: 24px;
        border-spacing: 0px;
      }
      .member-count {
        font-size: 18px;
      }
    `;
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
        this.requestUpdate();
      },
    });
    this.requestUpdate();
  }

  displayMemberStatus(memberLevel: number): string {
    const { MemberLevel } = MemberService;

    switch (memberLevel) {
      case MemberLevel.inactive:
        return "Inactive";
      case MemberLevel.student:
        return "Student";
      case MemberLevel.standard:
        return "Standard";
      case MemberLevel.premium:
        return "Premium";
      default:
        return "No member status found";
    }
  }
  render(): TemplateResult {
    return html`
      <card-element>
        <div class="member-container">
          <h1>Members</h1>
          <div class="member-count">
            <b>Member count: </b> ${this.memberCount}
          </div>
          <table>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Member Status</th>
              <th>Resources</th>
            </tr>
            ${this.members.map((x: MemberService.MemberResponse) => {
              return html`
                <tr>
                  <td class="name">${x.name}</td>
                  <td>${x.email}</td>
                  <td>${this.displayMemberStatus(x.memberLevel)}</td>
                  <td>No resources</td>
                </tr>
              `;
            })}
          </table>
        </div>
      </card-element>
    `;
  }
}
