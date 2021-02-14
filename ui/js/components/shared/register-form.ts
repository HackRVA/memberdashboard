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
import "@material/mwc-button";
import "@material/mwc-list";
import "@material/mwc-textfield";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-formfield";
import { TextField } from "@material/mwc-textfield/mwc-textfield";
// membership
import { showComponent } from "./../../function";
import { UserService } from "../../service/user.service";
import { defaultSnackbar } from "./default-snackbar";

@customElement("register-form")
export class RegisterForm extends LitElement {
  userService: UserService = new UserService();

  // form template
  emailFieldTemplate: TextField;
  passwordFieldTemplate: TextField;
  confirmPasswordFieldTemplate: TextField;

  static get styles(): CSSResult {
    return css`
      .register-container {
        height: 270px;
        max-width: 250px;
        background-color: #e1e1e1;
        padding: 24px;
      }
      mwc-formfield {
        display: block;
        margin-bottom: 16px;
      }
      .sign-in {
        float: left;
      }
      mwc-button {
        float: right;
        margin-bottom: 12px;
      }
    `;
  }

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot.querySelector("#email");
    this.passwordFieldTemplate = this.shadowRoot.querySelector("#password");
    this.confirmPasswordFieldTemplate = this.shadowRoot.querySelector(
      "#confirm-password"
    );
  }

  handleUserRegister(): void {
    const opts: UserService.RegisterRequest = {
      email: this.emailFieldTemplate?.value,
      password: this.passwordFieldTemplate?.value,
    };
    this.userService.registerUser(opts).subscribe({
      complete: () => this.displaySuccessMsg(),
      error: () => this.displayErrorMsg(),
    });
  }

  displaySuccessMsg(): void {
    showComponent("#success", this.shadowRoot);
  }

  displayErrorMsg(): void {
    showComponent("#error", this.shadowRoot);
  }

  displayInvalidMsg(): void {
    showComponent("#invalid", this.shadowRoot);
  }

  handleSubmit(): void {
    if (this.isPasswordIdentical() && this.isValid()) {
      this.handleUserRegister();
    } else {
      this.displayInvalidMsg();
    }
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.passwordFieldTemplate.validity.valid &&
      this.confirmPasswordFieldTemplate.validity.valid
    );
  }

  isPasswordIdentical(): boolean {
    return (
      this.passwordFieldTemplate.value ===
      this.confirmPasswordFieldTemplate.value
    );
  }

  fireSwitchEvent(): void {
    const switchToRegisterEvent = new CustomEvent("switch", {});
    this.dispatchEvent(switchToRegisterEvent);
  }

  goToLoginForm(): void {
    this.fireSwitchEvent();
  }

  render(): TemplateResult {
    return html`
      <div class="register-container">
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
        <mwc-formfield>
          <mwc-textfield
            id="confirm-password"
            required
            type="password"
            label="Retype password"
          ></mwc-textfield>
        </mwc-formfield>
        <mwc-button label="register" @click=${this.handleSubmit}></mwc-button>
        <mwc-button class="sign-in" @click=${this.goToLoginForm}>
          Sign in
        </mwc-button>
        ${defaultSnackbar("success", "success")}
        ${defaultSnackbar("invalid", "invalid")}
        ${defaultSnackbar("error", "error")}
      </div>
    `;
  }
}
