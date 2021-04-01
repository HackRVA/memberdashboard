// lit element
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
  property,
} from "lit-element";

// material
import { Checkbox } from "@material/mwc-checkbox";
import "@material/mwc-button";
import "@material/mwc-dialog";
import "@material/mwc-select";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-snackbar";
import "@material/mwc-textfield";
import "@material/mwc-checkbox";
import "@material/mwc-formfield";

// membership
import { MemberLevel, MemberResource, MemberResponse } from "./types";
import { showComponent } from "./../../function";
import { ResourceService, MemberService } from "../../service";
import { memberListStyles } from "./styles/member-list-styles";
import "../shared/rfid-modal";
import "./modals/add-member-to-resource-modal";
import "./modals/remove-member-from-resource-modal";
import "./modals/add-members-to-resource-modal";
import "../shared/card-element";

@customElement("member-list")
export class MemberList extends LitElement {
  @property({ type: Array })
  members: MemberResponse[] = [];

  @property({ type: Number })
  memberCount: number = 0;

  memberResources: Array<MemberResource> = [];
  email: string = "";

  memberEmails: string[] = [];

  membersCheckBoxTemplate: NodeListOf<Checkbox>;
  allMembersCheckBoxTemplate: Checkbox;

  memberService: MemberService = new MemberService();
  resourceService: ResourceService = new ResourceService();

  static get styles(): CSSResult[] {
    return [memberListStyles];
  }

  displayMemberStatus(memberLevel: MemberLevel): string {
    switch (memberLevel) {
      case MemberLevel.inactive:
        return "Inactive";
      case MemberLevel.student:
        return "Student";
      case MemberLevel.classic:
        return "Classic";
      case MemberLevel.standard:
        return "Standard";
      case MemberLevel.premium:
        return "Premium";
      default:
        return "No member status found";
    }
  }

  firstUpdated(): void {
    this.allMembersCheckBoxTemplate = this.shadowRoot.querySelector(
      "#all-members"
    );
  }

  updated(): void {
    this.membersCheckBoxTemplate = this.shadowRoot.querySelectorAll(
      "mwc-checkbox"
    );
  }

  getMembers(): void {
    this.memberService.getMembers().subscribe({
      next: (result: any) => {
        this.members = result as MemberResponse[];
        this.memberCount = this.members.length;
      },
      error: () => {
        console.error("unable to get members");
      },
    });
  }

  openAddMemberToResourceModal(email: string): void {
    this.email = email;
    this.requestUpdate();
    showComponent("#add-member-to-resource-modal", this.shadowRoot);
  }

  openAddMembersToResourceModal(): void {
    this.memberEmails = this.memberEmails.map((x: string) => x);
    this.requestUpdate();
    showComponent("#add-members-to-resource-modal", this.shadowRoot);
  }

  openRemoveMemberFromResourceModal(
    email: string,
    memberResources: Array<MemberResource>
  ): void {
    this.email = email;
    this.memberResources = memberResources;
    this.requestUpdate();
    showComponent("#remove-member-from-resource-modal", this.shadowRoot);
  }

  openRFIDModal(): void {
    showComponent("#rfid-modal", this.shadowRoot);
  }

  handleEmail(event: Event, email: string): void {
    this.allMembersCheckBoxTemplate.checked = false;

    const checked: boolean = (event.target as EventTarget & {
      checked: boolean;
    }).checked;

    if (checked) {
      this.memberEmails.push(email);
    } else {
      this.memberEmails = this.memberEmails.filter((x: string) => x !== email);
    }
  }

  handleAllEmails(event: Event): void {
    const checked: boolean = (event.target as EventTarget & {
      checked: boolean;
    }).checked;

    if (checked) {
      this.memberEmails = this.members.map((x: MemberResponse) => x.email);
      this.setAllCheckbox(true);
    } else {
      this.memberEmails = [];
      this.setAllCheckbox(false);
    }
  }

  setAllCheckbox(checked: boolean): void {
    this.membersCheckBoxTemplate.forEach((x: Checkbox) => {
      x.checked = checked;
    });
  }

  displayMembersTable(): TemplateResult {
    return html`
      ${this.members.map((x: MemberResponse) => {
        return html`
          <tr>
            <td class="name">
            <mwc-checkbox @change=${(event: Event) =>
              this.handleEmail(event, x.email)}></mwc-checkbox> 
            <span>${x.name}</span>
            </td>
            <td>${x.email}</td>
            <td>${this.displayMemberStatus(x.memberLevel)}</td>
            <td>
              <div class="horizontal-scrollbar">${this.displayMemberResources(
                x.resources
              )}</div>
              <div>
                <mwc-button
                  label="Add resource"
                  @click=${() => this.openAddMemberToResourceModal(x.email)}
                ></mwc-button>
                <mwc-button
                  class="remove"
                  label="Remove resource"
                  @click=${() =>
                    this.openRemoveMemberFromResourceModal(
                      x.email,
                      x.resources
                    )}
                ></mwc-button>
            </td>
            <td>
              ${x.rfid}
            </td>
          </tr>
        `;
      })}
    `;
  }

  displayMemberResources(resources: Array<MemberResource>): string {
    if (resources?.length > 0) {
      return resources.map((x: MemberResource) => x.name).join(", ");
    }
    return "No resources";
  }

  refreshMemberList(): void {
    this.getMembers();
    this.requestUpdate();
  }

  render(): TemplateResult {
    return html`
      <card-element>
        <div class="member-container">
          <div class="member-header">
            <h1>Members</h1>
            <span class="member-count">
              <b>Member count: </b> 
              ${this.memberCount} 
            </span>
            <mwc-button 
              class="rfid-button" 
              label="Assign rfid"
              unelevated 
              dense 
              @click=${this.openRFIDModal}> 
            </mvc-button>
          </div>
          <div class="all-members-action-container">
            <mwc-button
              class="add-resource-to-members"
              unelevated
              dense
              label="Add resource to members"
              @click=${this.openAddMembersToResourceModal}>
            </mwc-button>
            <mwc-formfield label="All members" class="all-members-checkbox">
              <mwc-checkbox id="all-members" @change=${
                this.handleAllEmails
              }></mwc-checkbox>
            </mwc-formfield>
          </div>
          <table>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Member Status</th>
              <th>Resources</th>
              <th>RFID</th>
            </tr>
            ${this.displayMembersTable()}
          </table>
        </div>
      </card-element>
      <add-members-to-resource-modal
        id="add-members-to-resource-modal"
        .emails=${this.memberEmails}
        @updated=${this.refreshMemberList}>
      </add-members-to-resource-modal>
      <add-member-to-resource-modal 
        id="add-member-to-resource-modal"
        .email=${this.email}
        @updated=${this.refreshMemberList}> 
      </add-member-to-resource-modal>
      <remove-member-from-resource-modal
        id="remove-member-from-resource-modal"
        .email=${this.email}
        .memberResources=${this.memberResources}
        @updated=${this.refreshMemberList}> 
      </remove-member-from-resource-modal>
      <rfid-modal
        id="rfid-modal"> 
      </rfid-modal>
    `;
  }
}
