import { showComponent } from "./../../function";
import { defaultSnackbar } from "./default-snackbar";
import {
  LitElement,
  html,
  TemplateResult,
  customElement,
  CSSResult,
  css,
} from "lit-element";
import { UserService } from "../../service/user.service";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-list/mwc-list-item";
import "@vaadin/vaadin-text-field/vaadin-text-field";
import "@vaadin/vaadin-text-field/vaadin-email-field";
import "@vaadin/vaadin-text-field/vaadin-password-field";
import "@vaadin/vaadin-form-layout";
import "@vaadin/vaadin-button";
import { EmailFieldElement } from "@vaadin/vaadin-text-field/vaadin-email-field";
import { TextFieldElement } from "@vaadin/vaadin-text-field/vaadin-text-field";
import { PasswordFieldElement } from "@vaadin/vaadin-text-field/vaadin-password-field";

@customElement("register-form")
export class RegisterForm extends LitElement {
  userService: UserService = new UserService();

  // form template
  usernameFieldTemplate: TextFieldElement;
  emailFieldTemplate: EmailFieldElement;
  passwordFieldTemplate: PasswordFieldElement;

  static get styles(): CSSResult {
    return css`
      vaadin-form-layout {
        max-width: 240px;
      }

      vaadin-text-field,
      vaadin-email-field,
      vaadin-password-field {
        margin-bottom: 12px;
      }

      vaadin-button {
        margin-top: 8px;
      }
    `;
  }

  firstUpdated(): void {
    this.usernameFieldTemplate = this.shadowRoot.querySelector(
      "#username-text-field"
    );
    this.emailFieldTemplate = this.shadowRoot.querySelector(
      "vaadin-email-field"
    );
    this.passwordFieldTemplate = this.shadowRoot.querySelector(
      "vaadin-password-field"
    );
  }

  handleUserRegister(): void {
    const opts: UserService.RegisterRequest = {
      username: this.usernameFieldTemplate?.value,
      password: this.passwordFieldTemplate?.value,
      email: this.emailFieldTemplate?.value,
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

  handleText(): void {
    if (this.isValid()) {
      this.handleUserRegister();
    } else {
      this.displayInvalidMsg();
    }
  }

  isValid(): boolean {
    return (
      this.usernameFieldTemplate?.validate() &&
      this.emailFieldTemplate?.validate() &&
      this.passwordFieldTemplate?.validate()
    );
  }

  render(): TemplateResult {
    return html`
      <div>
        <vaadin-form-layout>
          <vaadin-text-field
            id="username-text-field"
            required
            label="Username"
            placeholder="username"
            clear-button-visible
          ></vaadin-text-field>
          <vaadin-email-field
            required
            label="Email"
            placeholder="email"
            error-message="Please enter a valid email address"
            clear-button-visible
          ></vaadin-email-field>
          <vaadin-password-field
            required
            label="Password"
            reveal-button-hidden
            placeholder="password"
            clear-button-visible
          ></vaadin-password-field>
          <vaadin-button
            .disabled=${this.isValid()}
            theme="primary"
            @click=${this.handleText}
          >
            Register</vaadin-button
          >
        </vaadin-form-layout>
        ${defaultSnackbar("success", "success")}
        ${defaultSnackbar("invalid", "invalid")}
        ${defaultSnackbar("error", "error")}
      </div>
    `;
  }
}
