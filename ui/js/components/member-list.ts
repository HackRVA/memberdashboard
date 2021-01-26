import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  css,
  CSSResult,
  property,
} from "lit-element";
import { MemberService } from "../service/member.service";
import "./card-element";
import "@material/mwc-button";
import "@material/mwc-dialog";

@customElement("member-list")
export class MemberList extends LitElement {
  @property({ type: Array })
  members: MemberService.MemberResponse[] = [];

  @property({ type: Number })
  memberCount: number = 0;

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

  openMemberResourceModal(email: string): void {
    console.log("email", email);
    (this.shadowRoot?.querySelector("#memberResourceModal") as HTMLElement & {
      show: Function;
    }).show();
  }

  displayMembersTable(): TemplateResult {
    return html`
      ${this.members.map((x: MemberService.MemberResponse) => {
        return html`
          <tr>
            <td class="name">${x.name}</td>
            <td>${x.email}</td>
            <td>${this.displayMemberStatus(x.memberLevel)}</td>
            <td>
              <mwc-button
                label="Add resource"
                @click=${() => this.openMemberResourceModal(x.email)}
              ></mwc-button>
            </td>
          </tr>
        `;
      })}
    `;
  }

  displayMemberResourceModal(): TemplateResult {
    return html`
      <mwc-dialog id="memberResourceModal">
        <div>Add Resource</div>
        <mwc-button slot="primaryAction" dialogAction="discard">
          Discard
        </mwc-button>
        <mwc-button slot="secondaryAction" dialogAction="cancel">
          Cancel
        </mwc-button>
      </mwc-dialog>
    `;
  }

  displayMemberResources(
    resources: Array<MemberService.MemberResource>
  ): string {
    if (resources && resources.length > 0) {
      return resources.map((x) => x.name).join(", ");
    }
    return "No resources";
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
                  <td>${this.displayMemberResources(x.resources)}</td>
                </tr>
              `;
            })}
          </table>
        </div>

        ${this.displayMemberResourceModal()}
      </card-element>
    `;
  }
}
