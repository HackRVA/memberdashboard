// lit element
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
} from "lit-element";

// material
import "@material/mwc-snackbar";
import "@material/mwc-formfield";
import { TextField } from "@material/mwc-textfield/mwc-textfield";

// membership
import { UserService } from "../../service";
import { LoginRequest, Jwt } from "../user/types";
import { ToastMessage } from "../shared/types";
import { showComponent } from "../../function";
import { loginFormStyles } from "./styles";
import "../shared/toast-msg";

@customElement("login-form")
export class LoginForm extends LitElement {
  // form template
  emailFieldTemplate: TextField;
  passwordFieldTemplate: TextField;

  userService: UserService = new UserService();

  toastMsg: ToastMessage;

  static get styles(): CSSResult[] {
    return [loginFormStyles];
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
      this.displayToastMsg("Email and/or password invalid");
    }
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.passwordFieldTemplate.validity.valid
    );
  }

  handleUserLogin(): void {
    const opts: LoginRequest = {
      email: this.emailFieldTemplate?.value,
      password: this.passwordFieldTemplate?.value,
    };

    this.userService.login(opts).subscribe({
      next: (result: any) => {
        const { token } = result as Jwt;
        localStorage.setItem("jwt", token);
        window.location.reload();
      },
      error: () => {
        this.displayToastMsg("Email and/or password invalid");
      },
    });
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent("#toast-msg", this.shadowRoot);
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
        <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
      </div>
    `;
  }
}
