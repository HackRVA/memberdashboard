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
    `;
  }

  firstUpdated(): void {
    this.emailFieldTemplate = this.shadowRoot.querySelector("#email");
    this.passwordFieldTemplate = this.shadowRoot.querySelector("#password");
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
    if (this.isValid()) {
      this.handleUserRegister();
    } else {
      this.displayInvalidMsg();
    }
  }

  isValid(): boolean {
    return (
      this.emailFieldTemplate.validity.valid &&
      this.passwordFieldTemplate.validity.valid
    );
  }

  render(): TemplateResult {
    return html`
      <div>
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
        <mwc-button label="register" @click=${this.handleSubmit}></mwc-button>
        ${defaultSnackbar("success", "success")}
        ${defaultSnackbar("invalid", "invalid")}
        ${defaultSnackbar("error", "error")}
      </div>
    `;
  }
}
