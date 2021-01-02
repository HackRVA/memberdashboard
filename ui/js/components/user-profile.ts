import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import "./card-element";
import { UserService } from "../service/User";

@customElement("user-profile")
export class UserProfile extends LitElement {
  userService: UserService = new UserService();
  username: String = "";
  email: String = "";
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
  render(): TemplateResult {
    return html` <card-element>
      <profile-container>
      <profile-label>User Profile</profile-title>
      <username-label>${this.username}</username-label>
      <email-label>${this.email}</email-label>
       </profile-container>
    </card-element>`;
  }
}
