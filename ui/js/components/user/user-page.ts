import { defaultSnackbar } from "./../shared/default-snackbar";
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
import "../shared/card-element";
import { showComponent } from "../../function";
import { RFIDModal } from "../members/modals/rfid-modal";
import { MemberService } from "../../service/member.service";
import "@material/mwc-dialog";
import "@material/mwc-button";

@customElement("user-page")
export class UserPage extends LitElement {
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

      .email {
        font-size: 20px;
        line-height: 32px;
      }

      div {
        margin-bottom: 24px;
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
        const { email } = result as UserService.UserProfile;
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
    showComponent("#assignRFIDModal", this.shadowRoot);
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

  displaySuccessMessage(): void {
    showComponent("#success", this.shadowRoot);
  }

  displayErrorMessage(): void {
    showComponent("#error", this.shadowRoot);
  }

  assignUserToRFID(request: MemberService.AssignRFIDRequest): void {
    this.memberService.assignRFID(request).subscribe({
      complete: () => {
        this.displaySuccessMessage();
        this.requestUpdate();
      },
      error: () => {
        this.displayErrorMessage();
        this.requestUpdate();
      },
    });
  }

  render(): TemplateResult {
    return html`
    <div>
      <card-element>
        <h1> User <h1>
        <div class="center">
          <div> 
            <span class="email">${this.email} </span>
          </div>
          <div> 
            <mwc-button 
            class="rfid-button" 
            label="Assign rfid" 
            dense
            unelevated
            @click=${this.openRFIDModal}> 
            </mvc-button>
          </div> 
        </div>
        </card-element> 
        ${this.displayAddUpdateRFIDModal()}
        ${defaultSnackbar("error", "error")}
        ${defaultSnackbar("success", "success")}
    </div> 
    `;
  }
}
