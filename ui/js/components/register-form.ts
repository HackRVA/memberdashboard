import { LitElement, html, TemplateResult } from "lit-element";
import { USER_PROFILE_ACTOR_ADDRESS } from "../constants";
import ActorStore from "../actors/actorStore";
import { UserActor } from "../actors/user";
import "@material/mwc-textfield";
import "@material/mwc-button";
import "@material/mwc-list/mwc-list-item";

class RegisterForm extends LitElement {
  username: string = "";
  password: string = "";
  email: string = "";

  handleUsernameInput(e: KeyboardEvent): void {
    this.username = (e.target as HTMLInputElement).value;
  }

  handlePasswordInput(e: KeyboardEvent): void {
    this.password = (e.target as HTMLInputElement).value;
  }
  handleEmailInput(e: KeyboardEvent): void {
    this.email = (e.target as HTMLInputElement).value;
  }
  async handleUserRegister(): Promise<void> {
    const opts: UserActor.RegisterRequest = {
      username: this.username,
      password: this.password,
      email: this.email,
    };

    const userActor: any = ActorStore.lookup(USER_PROFILE_ACTOR_ADDRESS);
    const registerResponse: Promise<Boolean> = await userActor.message(
      UserActor.MessageTypes.RegisterUser,
      opts
    );

    const registerMessage = registerResponse
      ? "Success!"
      : "some kind of error registering user";

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

customElements.define("register-form", RegisterForm);
