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

// material
import "@material/mwc-dialog";
import "@material/mwc-button";
import "@material/mwc-textfield";

// membership
import { UserProfile } from "./types";
import { defaultSnackbar } from "./../shared/default-snackbar";
import { UserService } from "./../../service/user.service";
import { showComponent } from "../../function";
import { MemberService } from "../../service/member.service";
import "../shared/card-element";
import "../shared/rfid-modal";

@customElement("user-page")
export class UserPage extends LitElement {
  @property({ type: String })
  email: string = "";

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
        const { email } = result as UserProfile;
        this.email = email;
        this.requestUpdate();
      },
    });
  }

  openRFIDModal(): void {
    showComponent("#rfid-modal", this.shadowRoot);
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
        <rfid-modal 
          id="rfid-modal"
          .email=${this.email}
          >
        </rfid-modal>
        ${defaultSnackbar("error", "error")}
        ${defaultSnackbar("success", "success")}
    </div> 
    `;
  }
}
