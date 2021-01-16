import {
    LitElement,
    html,
    css,
    customElement,
    TemplateResult,
    CSSResult,
  } from "lit-element";
import { UserService } from "../service/User";
  
  @customElement("user-login-profile")
  export class UserLoginProfile extends LitElement {
    userService: UserService = new UserService();
    username: string = "";
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
          const { username, email } = result as UserService.UserProfile;
          this.username = username;
          this.email = email;
          this.requestUpdate();
        },
      });
    }

    handleLogout(): void {
        this.userService.logout().subscribe({
            next: (result) => {
                if ((result as { error: boolean; message: any }).error) {
                    console.log("error logging out");
                    return;
                }
                
                console.log("logging out");
                window.location.reload();
            },
        });
    }

    render(): TemplateResult {
      return html`
      <mwc-list-item>
        <mwc-icon slot="graphic">person</mwc-icon>
        ${this.username}</mwc-list-item
        >
        <mwc-list-item>${this.email}</mwc-list-item>
        <mwc-list-item @click=${this.handleLogout}>
        <mwc-button label="Logout"></mwc-button>
    </mwc-list-item>
      `;
    }
  }
  