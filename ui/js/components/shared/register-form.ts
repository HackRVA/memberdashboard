// lit element
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
} from "lit-element";

// material
import { TextField } from "@material/mwc-textfield/mwc-textfield";

// membership
import { showComponent } from "./../../function";
import { AuthService } from "../../service";
import { RegisterRequest } from "../shared/types";
import { ToastMessage } from "./types";
import { registerFormStyles } from "./styles/register-form-styles";
import "../shared/toast-msg";

@customElement("register-form")
export class RegisterForm extends LitElement {
  authService: AuthService = new AuthService();

  // form template
  emailFieldTemplate: TextField;
  passwordFieldTemplate: TextField;
  confirmPasswordFieldTemplate: TextField;

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [registerFormStyles];
  }

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot.querySelector("#email");
    this.passwordFieldTemplate = this.shadowRoot.querySelector("#password");
    this.confirmPasswordFieldTemplate = this.shadowRoot.querySelector(
      "#confirm-password"
    );
  }

  handleUserRegister(): void {
    const opts: RegisterRequest = {
      email: this.emailFieldTemplate?.value,
      password: this.passwordFieldTemplate?.value,
    };
    this.authService.registerUser(opts).subscribe({
      complete: () => this.displayToastMsg("Success"),
      error: () => this.displayToastMsg("Oops, something went wrong"),
    });
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
  }

  handleSubmit(): void {
    if (this.isPasswordIdentical() && this.isValid()) {
      this.handleUserRegister();
    } else {
      this.displayToastMsg(
        "Hrmmm, are you sure everything in the form is correct?"
      );
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

  render(): TemplateResult {
    return html`
      <mwc-formfield>
        <mwc-textfield
          size="30"
          id="email"
          required
          type="email"
          label="Email"
        ></mwc-textfield>
      </mwc-formfield>
      <mwc-formfield>
        <mwc-textfield
          size="30"
          id="password"
          required
          type="password"
          label="Password"
        ></mwc-textfield>
      </mwc-formfield>
      <mwc-formfield>
        <mwc-textfield
          size="30"
          id="confirm-password"
          required
          type="password"
          label="Retype password"
        ></mwc-textfield>
      </mwc-formfield>
      <mwc-button
        label="register"
        unelevated
        @click=${this.handleSubmit}
      ></mwc-button>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
