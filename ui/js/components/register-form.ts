import { LitElement, html, css, CSSResult, TemplateResult } from "lit-element";
import "@material/mwc-textfield";
import "@material/mwc-button";

interface RegisterRequest {
  Username: string;
  Password: string;
  Email: string;
}

console.log("help")
class RegisterForm extends LitElement {
  username: string = "";
  password: string = "";
  email: string = "";

  static get styles(): Array<CSSResult> {
    return [
      css`
        register-container {
          display: grid;
          width: 15%;
          grid-gap: 1em;
        }
      `,
    ];
  }

  handleUsernameInput(e: KeyboardEvent): void {
    this.username = (e.target as HTMLInputElement).value;
    console.log(this.username)
  }

  handlePasswordInput(e: KeyboardEvent): void {
    this.password = (e.target as HTMLInputElement).value;
    console.log(this.password)
  }
  handleEmailInput(e: KeyboardEvent): void {
    this.email = (e.target as HTMLInputElement).value;
    console.log(this.email)
  }
  handleUserRegister(): void {
      console.log("btn-cliecked")
    const opts: RegisterRequest = {
      Username: this.username,
      Password: this.password,
      Email: this.email
    };
    fetch("/register", {
      method: "POST",
      body: JSON.stringify(opts),
    })
      .then((response) => response.json())
      .then(console.log);
  }
  render(): TemplateResult {
    return html`<register-container>
      <mwc-textfield label="Username" @change=${this.handleUsernameInput}></mwc-textfield>
      <mwc-textfield type="email" label="Email" @change=${this.handleEmailInput}></mwc-textfield>
      <mwc-textfield type="password" label="Password" @change=${this.handlePasswordInput}></mwc-textfield>
      <mwc-button label="register" @click=${this.handleUserRegister}></mwc-button>
    </register-container>`;
  }
}

customElements.define("register-form", RegisterForm);
