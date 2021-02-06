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

@customElement("user-profile")
export class UserProfile extends LitElement {
  @property({ type: String })
  username: string = "";

  @property({ type: String })
  email: string = "";

  userService: UserService = new UserService();
  static get styles(): CSSResult {
    return css`
      .user-profile-container {
        padding: 32px;
        margin-bottom: 24px;
      }

      .username,
      .email {
        font-size: 20px;
        line-height: 32px;
      }
    `;
  }

  render(): TemplateResult {
    return html`
      <div class="user-profile-container">
        <div class="username">${this.username}</div>
        <div class="email">${this.email}</div>
      </div>
    `;
  }
}
