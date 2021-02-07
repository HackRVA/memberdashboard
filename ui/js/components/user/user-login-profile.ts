import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import { UserService } from "../../service/user.service";
import "@material/mwc-button";

@customElement("user-login-profile")
export class UserLoginProfile extends LitElement {
  userService: UserService = new UserService();
  email: string = "";
  static get styles(): CSSResult {
    return css``;
  }

  firstUpdated(): void {
    this.handleGetUserProfile();
  }

  handleGetUserProfile(): void {
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
      <mwc-list-item>
        <mwc-icon slot="graphic">person</mwc-icon>
        ${this.email}</mwc-list-item
      >
      <mwc-list-item @click=${this.handleLogout}>
        <mwc-button label="Logout"></mwc-button>
      </mwc-list-item>
    `;
  }
}
