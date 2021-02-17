// lit element
import {
  LitElement,
  html,
  css,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import { UserService } from "../../service";
import "../shared/card-element";
import "../shared/register-form";
import { UserProfile } from "../user/types";

@customElement("home-page")
export class HomePage extends LitElement {
  userService: UserService = new UserService();
  isUserLogin: boolean = false;
  isRegister: boolean = false;

  static get styles(): CSSResult {
    return css`
      .login-container {
        display: grid;
        justify-content: center;
        padding: 36px;
      }

      .center {
        text-align: center;
      }
    `;
  }

  firstUpdated(): void {
    this.checkUserLogin();
  }

  checkUserLogin(): void {
    this.userService.getUser().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        const { email } = result as UserProfile;
        this.isUserLogin = !!email;
        this.requestUpdate();
      },
    });
  }

  displayRegisterLoginForm(): TemplateResult {
    if (this.isRegister) {
      return html`<register-form @switch=${this.handleSwitch} />`;
    } else {
      return html`<login-form @switch=${this.handleSwitch} />`;
    }
  }

  handleSwitch(): void {
    this.isRegister = !this.isRegister;
    this.requestUpdate();
  }

  displayHomePage(): TemplateResult {
    if (this.isUserLogin) {
      return html` <h1>Home</h1> `;
    } else {
      return html`
        <div>
          <h1>Home</h1>
          <div class="login-container">${this.displayRegisterLoginForm()}</div>
        </div>
      `;
    }
  }

  render(): TemplateResult {
    return html` <card-element> ${this.displayHomePage()} </card-element> `;
  }
}
