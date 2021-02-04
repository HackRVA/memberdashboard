import { ResourceService } from "../../service/resource.service";
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  css,
  CSSResult,
  property,
} from "lit-element";
import { MemberService } from "../../service/member.service";
import "../shared/card-element";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-select";
import "@material/mwc-list/mwc-list-item";

@customElement("member-list")
export class MemberList extends LitElement {
  @property({ type: Array })
  members: MemberService.MemberResponse[] = [];

  @property({ type: Number })
  memberCount: number = 0;

  resources: ResourceService.ResourceResponse[] = [];

  // form variables for adding/removing a resource to a member
  email: string = "";
  newResourceId: string = "";

  memberResources: Array<MemberService.MemberResource> = [];
  memberService: MemberService = new MemberService();
  resourceService: ResourceService = new ResourceService();

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
        font-size: 20px;
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
      .remove {
        --mdc-theme-primary: #e9437a;
      }
      .horizontal-scrollbar {
        overflow: auto;
        max-width: 320px;
        white-space: nowrap;
      }
    `;
  }

  firstUpdated(): void {
    this.getResources();
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

  getResources(): void {
    this.resourceService.getResources().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          // this.onLoginComplete("Some error logging in");
          console.error("some error getting resources");
        } else {
          this.resources = result as ResourceService.ResourceResponse[];
          this.requestUpdate();
        }
      },
    });
  }

  getMembers(): void {
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

  openAddMemberResourceModal(email: string): void {
    this.email = email;
    this.requestUpdate();

    (this.shadowRoot?.querySelector(
      "#addMemberResourceModal"
    ) as HTMLElement & {
      show: Function;
    }).show();
  }

  openRemoveMemberResourceModal(
    email: string,
    memberResources: Array<MemberService.MemberResource>
  ): void {
    this.email = email;
    this.memberResources = memberResources;
    this.requestUpdate();

    (this.shadowRoot?.querySelector(
      "#removeMemberResourceModal"
    ) as HTMLElement & {
      show: Function;
    }).show();
  }

  handleResourceChange(e: Event): void {
    this.newResourceId = (e.target as EventTarget & { value: string }).value;
  }

  handleSubmitAddMemberResource(): void {
    const request: ResourceService.AddMemberResourceRequest = {
      email: this.email,
      resourceID: this.newResourceId,
    };
    this.emptyFormValues();
    this.addMemberResource(request);
  }

  handleSubmitRemoveMemberResource(): void {
    const request: ResourceService.RemoveMemberResourceRequest = {
      email: this.email,
      resourceID: this.newResourceId,
    };
    this.emptyFormValues();
    this.removeMemberResource(request);
  }

  removeMemberResource(
    request: ResourceService.RemoveMemberResourceRequest
  ): void {
    this.resourceService.removeMemberResource(request).subscribe({
      complete: () => {
        this.getMembers();
        this.requestUpdate();
      },
    });
  }

  addMemberResource(request: ResourceService.AddMemberResourceRequest): void {
    this.resourceService.addMemberResource(request).subscribe({
      complete: () => {
        this.getMembers();
        this.requestUpdate();
      },
    });
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
              <div class="horizontal-scrollbar">${this.displayMemberResources(
                x.resources
              )}</div>
              <div>
                <mwc-button
                  label="Add resource"
                  @click=${() => this.openAddMemberResourceModal(x.email)}
                ></mwc-button>
                <mwc-button
                  class="remove"
                  label="Remove resource"
                  @click=${() =>
                    this.openRemoveMemberResourceModal(x.email, x.resources)}
                ></mwc-button>
            </td>
          </tr>
        `;
      })}
    `;
  }

  displayRemoveMemberResourceModal(): TemplateResult {
    return html`
      <mwc-dialog id="removeMemberResourceModal">
        <div>Remove Resource</div>
        <mwc-textfield
          label="email"
          helper="Can't edit email"
          readonly
          value=${this.email}
        ></mwc-textfield>
        <mwc-select label="Resources" @change=${this.handleResourceChange}>
          ${this.memberResources.map((x: MemberService.MemberResource) => {
            return html`
              <mwc-list-item value=${x.resourceID}> ${x.name} </mwc-list-item>
            `;
          })}
        </mwc-select>
        <mwc-button
          slot="primaryAction"
          dialogAction="ok"
          @click=${this.handleSubmitRemoveMemberResource}
        >
          Submit
        </mwc-button>
        <mwc-button
          slot="secondaryAction"
          dialogAction="cancel"
          @click=${this.emptyFormValuesOnClosed}
        >
          Cancel
        </mwc-button>
      </mwc-dialog>
    `;
  }

  displayAddMemberResourceModal(): TemplateResult {
    return html`
      <mwc-dialog id="addMemberResourceModal">
        <div>Add Resource</div>
        <mwc-textfield
          label="email"
          helper="Can't edit email"
          readonly
          value=${this.email}
        ></mwc-textfield>
        <mwc-select label="Resources" @change=${this.handleResourceChange}>
          ${this.resources.map((x: ResourceService.ResourceResponse) => {
            return html`
              <mwc-list-item value=${x.id}> ${x.name} </mwc-list-item>
            `;
          })}
        </mwc-select>
        <mwc-button
          slot="primaryAction"
          dialogAction="ok"
          @click=${this.handleSubmitAddMemberResource}
        >
          Submit
        </mwc-button>
        <mwc-button
          slot="secondaryAction"
          dialogAction="cancel"
          @click=${this.emptyFormValuesOnClosed}
        >
          Cancel
        </mwc-button>
      </mwc-dialog>
    `;
  }

  emptyFormValuesOnClosed(): void {
    this.emptyFormValues();
    this.requestUpdate();
  }
  emptyFormValues(): void {
    this.email = "";
    this.newResourceId = "";
  }

  displayMemberResources(
    resources: Array<MemberService.MemberResource>
  ): string {
    if (resources?.length > 0) {
      return resources
        .map((x: MemberService.MemberResource) => x.name)
        .join(", ");
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
            ${this.displayMembersTable()}
          </table>
        </div>

        ${this.displayAddMemberResourceModal()}
        ${this.displayRemoveMemberResourceModal()}
      </card-element>
    `;
  }
}
