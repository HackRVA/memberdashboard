import { LitElement, html, TemplateResult } from "lit-element";
import { USER_PROFILE_ACTOR_ADDRESS } from "../constants";
import ActorStore from "../actors/actorStore";
import { UserActor } from "../actors/user";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-snackbar";
import "@material/mwc-list/mwc-list-item";

class LoginForm extends LitElement {
  username: string = "";
  password: string = "";
  loginMessage: string = "";

  handleUsernameInput(e: KeyboardEvent): void {
    this.username = (e.target as HTMLInputElement).value;
  }

  handlePasswordInput(e: KeyboardEvent): void {
    this.password = (e.target as HTMLInputElement).value;
  }

  async handleUserLogin(): Promise<void> {
    const opts: UserActor.LoginRequest = {
      username: this.username,
      password: this.password,
    };

    const userActor: any = ActorStore.lookup(USER_PROFILE_ACTOR_ADDRESS);
    const loginResponse: Promise<Boolean> = await userActor.message(
      UserActor.MessageTypes.Login,
      opts
    );
    const loginMessage = loginResponse
      ? "Success!"
      : "some kind of error logging in";

    const event = new CustomEvent("control-changed", {
      detail: loginMessage,
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
          type="password"
          label="Password"
          @change=${this.handlePasswordInput}
        ></mwc-textfield>
      </mwc-list-item>
      <mwc-list-item>
        <mwc-button label="login" @click=${this.handleUserLogin}></mwc-button>
      </mwc-list-item>
    `;
  }
}

customElements.define("login-form", LoginForm);
