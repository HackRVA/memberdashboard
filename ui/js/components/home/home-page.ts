import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";
import { UserService } from "../../service/user.service";
import "../shared/card-element";
import "../shared/register-form";

@customElement("home-page")
export class HomePage extends LitElement {
  email: string = "";
  userService: UserService = new UserService();

  static get styles(): CSSResult {
    return css`
      login-container {
        display: grid;
        justify-content: center;
        padding: 24px;
      }

      .center {
        text-align: center;
      }
    `;
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

  displayHomePage(): TemplateResult {
    if (this.email) {
      return html` <h1>Home</h1> `;
    } else {
      return html`
        <card-element>
          <h1>Home</h1>
          <login-container>
            <register-form />
          </login-container>
        </card-element>
      `;
    }
  }

  render(): TemplateResult {
    return html` ${this.displayHomePage()} `;
  }
}
