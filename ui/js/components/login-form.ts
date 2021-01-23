import { LitElement, html, TemplateResult, customElement } from "lit-element";
import { UserService } from "../service/user.service";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-snackbar";
import "@material/mwc-list/mwc-list-item";

@customElement("login-form")
export class LoginForm extends LitElement {
  username: string = "";
  password: string = "";
  userService: UserService = new UserService();

  onLoginComplete(snackbarNotification: string): void {
    const event = new CustomEvent("control-changed", {
      detail: snackbarNotification,
    });
    this.dispatchEvent(event);
    // window.location.reload();
  }

  handleUsernameInput(e: KeyboardEvent): void {
    this.username = (e.target as HTMLInputElement).value;
  }

  handlePasswordInput(e: KeyboardEvent): void {
    this.password = (e.target as HTMLInputElement).value;
  }

  handleUserLogin(): void {
    const opts: UserService.LoginRequest = {
      username: this.username,
      password: this.password,
    };
    this.userService
      .login(opts)
      .subscribe((response: UserService.ILoginResponse) => {
        if (response.token) {
          localStorage.setItem("jwt", response.token);
          window.location.reload();
        }
      });
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
          type="password"
          label="Password"
          @change=${this.handlePasswordInput}
        ></mwc-textfield>
      </mwc-list-item>
      <mwc-list-item @click=${this.handleUserLogin}>
        <mwc-button label="login"></mwc-button>
      </mwc-list-item>
    `;
  }
}
