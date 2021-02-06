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

@customElement("user-page")
export class UserPage extends LitElement {
  @property({ type: String })
  username: string = "";

  @property({ type: String })
  email: string = "";

  userService: UserService = new UserService();

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

  render(): TemplateResult {
    return html` 
    <card-element class="center">
      <h1> User <h1>
      <user-profile .username=${this.username} .email=${this.email} />
    </card-element> 
    `;
  }
}
