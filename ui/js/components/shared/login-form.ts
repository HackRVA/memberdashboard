import { LitElement, html, TemplateResult, customElement } from "lit-element";
import { UserService } from "../../service/user.service";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-snackbar";
import "@material/mwc-list/mwc-list-item";

@customElement("login-form")
export class LoginForm extends LitElement {
  email: string = "";
  password: string = "";
  userService: UserService = new UserService();

  handleEmailInput(e: KeyboardEvent): void {
    this.email = (e.target as HTMLInputElement).value;
  }

  handlePasswordInput(e: KeyboardEvent): void {
    this.password = (e.target as HTMLInputElement).value;
  }

  handleUserLogin(): void {
    const opts: UserService.LoginRequest = {
      email: this.email,
      password: this.password,
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
      <mwc-list-item>
        <mwc-textfield
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
      <mwc-list-item @click=${this.handleUserLogin}>
        <mwc-button label="login"></mwc-button>
      </mwc-list-item>
    `;
  }
}
