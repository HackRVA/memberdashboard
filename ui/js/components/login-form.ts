import { LitElement, html, css, CSSResult, TemplateResult } from "lit-element";
import "@material/mwc-textfield";
import "@material/mwc-button";

interface LoginRequest {
  Username: string;
  Password: string;
}

class LoginForm extends LitElement {
  username: string = "";
  password: string = "";

  static get styles(): Array<CSSResult> {
    return [
      css`
        login-container {
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

  handleUserLogin(): void {
      console.log("btn-cliecked")
    const opts: LoginRequest = {
      Username: this.username,
      Password: this.password,
    };
    fetch("/signin", {
      method: "POST",
      body: JSON.stringify(opts),
    })
      .then((response) => response.json())
      .then(console.log);
  }
  render(): TemplateResult {
    return html`<login-container>
      <mwc-textfield label="Username" @change=${this.handleUsernameInput}></mwc-textfield>
      <mwc-textfield type="password" label="Password" @change=${this.handlePasswordInput}></mwc-textfield>
      <mwc-button label="login" @click=${this.handleUserLogin}></mwc-button>
    </login-container>`;
  }
}

customElements.define("login-form", LoginForm);
