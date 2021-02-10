import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";
import { UserService } from "../../service/user.service";
import "@material/mwc-button";

@customElement("user-login-profile")
export class UserLoginProfile extends LitElement {
  @property({ type: String })
  email: string;

  userService: UserService = new UserService();
  static get styles(): CSSResult {
    return css`
      .logout {
        float: right;
      }
    `;
  }

  handleLogout(): void {
    this.userService.logout().subscribe({
      next: (response: null) => {
        localStorage.removeItem("jwt");
        window.location.reload();
      },
    });
  }

  render(): TemplateResult {
    return html`
      <div>
        <mwc-list-item>
          Signed in as <strong>${this.email} </strong></mwc-list-item
        >
        <mwc-list-item class="logout" @click=${this.handleLogout}>
          <mwc-button label="Logout"></mwc-button>
        </mwc-list-item>
      </div>
    `;
  }
}
