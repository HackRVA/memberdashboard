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
      mwc-formfield {
        display: block;
        margin-bottom: 16px;
      }
      mwc-button {
        float: right;
        margin-bottom: 12px;
      }
      .login-form {
        padding: 20px;
      }
    `;
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

  render(): TemplateResult {
    return html`
      <div class="login-form">
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
        <mwc-button label="login" @click=${this.handleSubmit}></mwc-button>
      </div>
    `;
  }
}
