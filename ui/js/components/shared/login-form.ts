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

// vaadin
import "@vaadin/vaadin-login/vaadin-login-form";
import { LoginFormElement } from "@vaadin/vaadin-login/vaadin-login-form";
import { LoginI18n } from "@vaadin/vaadin-login/src/interfaces";

// membership
import { UserService } from "../../service/user.service";

@customElement("login-form")
export class LoginForm extends LitElement {
  userService: UserService = new UserService();
  loginFormTemplate: LoginFormElement;

  static get styles(): CSSResult {
    return css`
      vaadin-button {
        background-color: #6200ee;
      }
    `;
  }

  firstUpdated(): void {
    this.loginFormTemplate = this.shadowRoot?.querySelector(
      "vaadin-login-form"
    );

    this.loginFormTemplate.i18n = this.updateI18n();
  }

  updateI18n(): LoginI18n {
    const newLoginil8n = {
      form: {
        title: "Welcome back",
        username: "Email address",
        password: "Password",
        submit: "Log in",
        forgotPassword: "", // eventually add this back in
      },
    };

    const i18n: LoginI18n = Object.assign(
      {},
      this.loginFormTemplate.i18n,
      newLoginil8n
    );

    return i18n;
  }

  handleUserLogin(event: CustomEvent): void {
    const opts: UserService.LoginRequest = {
      username: event.detail.username,
      password: event.detail.password,
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
      <vaadin-login-form @login=${this.handleUserLogin}> </vaadin-login-form>
    `;
  }
}
