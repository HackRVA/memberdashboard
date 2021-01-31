import { LitElement, html, TemplateResult, customElement } from "lit-element";
import { UserService } from "../../service/user.service";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-list/mwc-list-item";

@customElement("register-form")
export class RegisterForm extends LitElement {
  username: string = "";
  password: string = "";
  email: string = "";
  userService: UserService = new UserService();

  handleUsernameInput(e: KeyboardEvent): void {
    this.username = (e.target as HTMLInputElement).value;
  }

  handlePasswordInput(e: KeyboardEvent): void {
    this.password = (e.target as HTMLInputElement).value;
  }
  handleEmailInput(e: KeyboardEvent): void {
    this.email = (e.target as HTMLInputElement).value;
  }
  handleUserRegister(): void {
    const opts: UserService.RegisterRequest = {
      username: this.username,
      password: this.password,
      email: this.email,
    };
    this.userService.registerUser(opts).subscribe({
      next: (result) => {
        if ((result as { error: boolean; message: any }).error) {
          this.onRegisterComplete("Some error logging in");
        }
      },
      complete: () => this.onRegisterComplete("Success!"),
    });
  }

  onRegisterComplete(registerMessage: String) {
    const event = new CustomEvent("control-changed", {
      detail: registerMessage,
    });
    this.dispatchEvent(event);
  }

  render(): TemplateResult {
    return html`
      <mwc-list-item>
        <mwc-textfield
          label="Username"
          @change=${this.handleUsernameInput}
        ></mwc-textfield>
      </mwc-list-item>
      <mwc-list-item>
        <mwc-textfield
          type="email"
          label="Email"
          @change=${this.handleEmailInput}
        ></mwc-textfield>
      </mwc-list-item>
      <mwc-list-item>
        <mwc-textfield
          type="password"
          label="Password"
          @change=${this.handlePasswordInput}
        ></mwc-textfield>
      </mwc-list-item>
      <mwc-list-item @click=${this.handleUserRegister}>
        <mwc-button label="register"></mwc-button>
      </mwc-list-item>
    `;
  }
}
