// lit element
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
} from "lit-element";

// material
import "@material/mwc-button";
import "@material/mwc-list";
import "@material/mwc-textfield";
import "@material/mwc-list/mwc-list-item";
import "@material/snackbar";
import "@material/mwc-formfield";
import { TextField } from "@material/mwc-textfield/mwc-textfield";

// membership
import { showComponent } from "./../../function";
import { UserService } from "../../service";
import { RegisterRequest } from "../user/types";
import { ToastMessage } from "./types";
import { registerFormStyles } from "./styles/register-form-styles";
import "../shared/toast-msg";

@customElement("register-form")
export class RegisterForm extends LitElement {
  userService: UserService = new UserService();

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
    this.userService.registerUser(opts).subscribe({
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
        <mwc-button class="sign-in" @click=${this.goToLoginForm}>
          Sign in
        </mwc-button>
        <mwc-button label="register" @click=${this.handleSubmit}></mwc-button>
        <toast-msg
          id="toast-msg"
          .toastMsg=${this.toastMsg}> 
        </toast-msg>
      </card-element>
      </div>
    `;
  }
}
