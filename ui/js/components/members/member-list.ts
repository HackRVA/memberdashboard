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

// membership
import { MemberLevel, MemberResource, MemberResponse } from "./types";
import { showComponent } from "./../../function";
import { ResourceService, MemberService, PaymentService } from "../../service";
import { memberListStyles } from "./styles/member-list-styles";
import "../shared/rfid-modal";
import "./modals/add-member-to-resource-modal";
import "./modals/remove-member-from-resource-modal";
import "./modals/add-members-to-resource-modal";
import "../shared/card-element";
import { ToastMessage } from "../shared/types";

@customElement("member-list")
export class MemberList extends LitElement {
  @property({ type: Array })
  members: MemberResponse[] = [];

  @property({ type: Number })
  memberCount: number = 0;

  memberResources: Array<MemberResource> = [];
  email: string = "";

  memberEmails: string[] = [];
  toastMsg: ToastMessage;

  membersCheckBoxTemplate: NodeListOf<Checkbox>;
  allMembersCheckBoxTemplate: Checkbox;

  memberService: MemberService = new MemberService();
  resourceService: ResourceService = new ResourceService();
  paymentsService: PaymentService = new PaymentService();

  static get styles(): CSSResult[] {
    return [memberListStyles];
  }

  displayMemberStatus(memberLevel: MemberLevel): string {
    switch (memberLevel) {
      case MemberLevel.inactive:
        return "Inactive";
      case MemberLevel.credited:
        return "Credited";
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

  refreshMembersPayments(): void {
    this.paymentsService.refreshPayments().subscribe({
      complete: () => {
        this.refreshMemberList();
        this.displayToastMsg("Success");
      },
    });
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  render(): TemplateResult {
    return html`
      <card-element>
        <div class="member-container">
          <div class="member-header">
            <h3 class="member-count">
              Number of active members: ${this.memberCount} 
            </h3>
            <div class="buttons-container">
              <mwc-button 
                class="refresh-members-list" 
                unelevated 
                dense 
                label="Refresh member list"
                @click=${this.refreshMembersPayments}> 
              </mwc-button>
              <mwc-button 
                class="rfid-button" 
                label="Assign rfid"
                unelevated 
                dense 
                @click=${this.openRFIDModal}> 
              </mvc-button>
            </div>
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
        id="rfid-modal"
        @updated=${this.refreshMemberList}> 
      </rfid-modal>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
