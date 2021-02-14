// lit element
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
  css,
} from "lit-element";

// material
import "@material/mwc-snackbar";
import "@material/mwc-formfield";
import { TextField } from "@material/mwc-textfield/mwc-textfield";

// membership
import { UserService } from "../../service/user.service";

@customElement("login-form")
export class LoginForm extends LitElement {
  // form template
  emailFieldTemplate: TextField;
  passwordFieldTemplate: TextField;

  userService: UserService = new UserService();

  static get styles(): CSSResult {
    return css`
      .login-container {
        height: 270px;
        max-width: 250px;
        background-color: #e1e1e1;
        padding: 24px;
      }
      mwc-formfield {
        display: block;
        margin-bottom: 16px;
      }
      .mwc-button {
        margin-bottom: 12px;
      }
      .register {
        float: left;
      }
      .login {
        float: right;
      }
    `;
  }

  fireSwitchEvent(): void {
    const switchToRegisterEvent = new CustomEvent("switch", {});
    this.dispatchEvent(switchToRegisterEvent);
  }

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot?.querySelector("#email");
    this.passwordFieldTemplate = this.shadowRoot?.querySelector("#password");
  }

  handleSubmit(): void {
    if (this.isValid()) {
      this.handleUserLogin();
    } else {
      console.error("invalid");
    }
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.passwordFieldTemplate.validity.valid
    );
  }

  handleUserLogin(): void {
    const opts: UserService.LoginRequest = {
      email: this.emailFieldTemplate?.value,
      password: this.passwordFieldTemplate?.value,
    };
    this.userService.login(opts).subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        const { token } = result as UserService.Jwt;
        localStorage.setItem("jwt", token);
        window.location.reload();
      },
    });
  }

  goToRegisterForm(): void {
    this.fireSwitchEvent();
  }

  render(): TemplateResult {
    return html`
      <div class="login-container">
        <mwc-formfield>
          <mwc-textfield
            id="email"
            required
            type="email"
            label="Email"
          ></mwc-textfield>
        </mwc-formfield>
        <mwc-formfield>
          <mwc-textfield
            id="password"
            required
            type="password"
            label="Password"
          ></mwc-textfield>
        </mwc-formfield>
        <mwc-button
          label="Sign up"
          class="register"
          @click=${this.goToRegisterForm}
        ></mwc-button>
        <mwc-button
          label="login"
          @click=${this.handleSubmit}
          class="login"
        ></mwc-button>
      </div>
    `;
  }
}
