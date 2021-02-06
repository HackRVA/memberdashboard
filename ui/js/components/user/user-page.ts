import { UserService } from "./../../service/user.service";
import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";
import "./user-profile";
import "../shared/card-element";
import { openComponent } from "../../function";
import { RFIDModal } from "../members/modals/rfid-modal";
import { MemberService } from "../../service/member.service";
import "@material/mwc-dialog";
import "@material/mwc-button";

@customElement("user-page")
export class UserPage extends LitElement {
  @property({ type: String })
  username: string = "";

  @property({ type: String })
  email: string = "";

  newRFID: string = "";

  userService: UserService = new UserService();
  memberService: MemberService = new MemberService();

  static get styles(): CSSResult {
    return css`
      .center {
        text-align: center;
      }
    `;
  }

  firstUpdated(): void {
    this.userService.getUser().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        const { username, email } = result as UserService.UserProfile;
        this.username = username;
        this.email = email;
        this.requestUpdate();
      },
    });
  }

  displayAddUpdateRFIDModal(): TemplateResult {
    const modalData: MemberService.RFIDModalData = {
      email: this.email,
      rfid: this.newRFID,
      handleEmailChange: this.handleEmailChange,
      handleRFIDChange: this.handleRFIDChange,
      handleSubmitForAssigningMemberToRFID: this
        .handleSubmitForAssigningUserToRFID,
      emptyFormValuesOnClosed: this.emptyFormValuesOnClosed,
    };
    return RFIDModal(modalData);
  }

  openRFIDModal(): void {
    openComponent("#assignRFIDModal", this.shadowRoot);
  }

  emptyFormValues(): void {
    this.newRFID = "";
  }

  emptyFormValuesOnClosed(): void {
    this.emptyFormValues();
    this.requestUpdate();
  }

  handleEmailChange(e: Event): void {
    this.email = (e.target as EventTarget & { value: string }).value;
  }

  handleRFIDChange(e: Event): void {
    this.newRFID = (e.target as EventTarget & { value: string }).value;
  }

  handleSubmitForAssigningUserToRFID(): void {
    const request: MemberService.AssignRFIDRequest = {
      email: this.email.trim(),
      rfid: this.newRFID,
    };
    this.emptyFormValues();
    this.assignUserToRFID(request);
  }

  assignUserToRFID(request: MemberService.AssignRFIDRequest): void {
    this.memberService.assignRFID(request).subscribe();
  }

  render(): TemplateResult {
    return html`
    <div>
      <card-element class="center">
        <div> 
          <h1> User <h1>
          <mwc-button 
          class="rfid-button" 
          label="Assign rfid" 
          dense 
          @click=${this.openRFIDModal}> </mvc-button>
        </div>
        <user-profile .username=${this.username} .email=${this.email} />
        </card-element> 
        ${this.displayAddUpdateRFIDModal()}
    </div> 
    `;
  }
}
