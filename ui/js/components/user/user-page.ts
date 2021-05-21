// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";

// membership
import { UserProfile } from "./types";
import { showComponent } from "../../function";
import { UserService } from "../../service";
import { userPageStyles } from "./styles";
import "../shared/card-element";
import "../shared/rfid-modal";

@customElement("user-page")
export class UserPage extends LitElement {
  @property({ type: String })
  email: string = "";

  userService: UserService = new UserService();

  static get styles(): CSSResult[] {
    return [userPageStyles];
  }

  firstUpdated(): void {
    this.userService.getUser().subscribe({
      next: (result: UserProfile) => {
        const { email } = result;
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
    <card-element>
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
    <rfid-modal id="rfid-modal" .email=${this.email}>
    </rfid-modal>
    `;
  }
}
