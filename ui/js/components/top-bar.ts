import { LitElement, html, css, CSSResult, TemplateResult } from "lit-element";
import { UserService } from "../service/User";
import "./login-form";
import "./register-form";
import "./member-stats";
import "@material/mwc-top-app-bar-fixed";
import "@material/mwc-icon-button";
import "@material/mwc-menu";
import "@material/mwc-list/mwc-list-item";
import "@material/mwc-icon";

class TopBar extends LitElement {
  showRegister: Boolean = false;
  showUserProfile: Boolean = false;
  snackMessage: String = "";
  username: String = "";
  email: String = "";
  userService: UserService = new UserService();

  static get styles(): CSSResult {
    return css`
      login-container {
        display: grid;
        justify-content: center;
        margin-top: 10vh;
      }
    `;
  }

  updated(): void {
    if (this.showUserProfile) return;

    this.userService.getUser().subscribe({
      next: (result: any) => {
        if ((result as { error: boolean; message: any }).error) {
          return console.error(
            (result as { error: boolean; message: any }).message
          );
        }
        const { username, email } = result as UserService.UserProfile;
        this.username = username;
        this.email = email;
        this.showUserProfile = true;
        this.requestUpdate();
      },
    });
  }

  handleLogout(): void {
    this.userService.logout().subscribe({
      next: (result) => {
        if ((result as { error: boolean; message: any }).error) {
          this.writeToSnackbar("error logging out");
          return;
        }

        this.writeToSnackbar("logging out");
      },
    });
  }

  handleSnackbarMsg(evt: Event & { detail: String }): void {
    this.writeToSnackbar(evt.detail);
  }

  writeToSnackbar(message: String): void {
    const snackbar:
      | (HTMLElement & { show: Function })
      | null
      | undefined = this.shadowRoot?.querySelector("#loginMessage");
    if (!snackbar) return console.error("no snackbar");

    this.snackMessage = message;

    this.requestUpdate();
    snackbar.show();
  }

  handleRegisterBtn(): void {
    this.showRegister = !this.showRegister;
    this.requestUpdate();
  }

  handleProfileClick(): void {
    const profileBtn:
      | HTMLElement
      | null
      | undefined = this.shadowRoot?.querySelector("#profileBtn");
    const menu:
      | (HTMLElement & { anchor: HTMLElement; show: Function })
      | null
      | undefined = this.shadowRoot?.querySelector("#menu");

    if (!profileBtn) return console.error("profile btn doesn't exist");
    if (!menu) return console.error("menu element doesn't exist");

    menu.anchor = profileBtn;
    menu.show();
  }

  render(): TemplateResult {
    const login = html`<login-form
        @control-changed="${this.handleSnackbarMsg}"
      ></login-form>
      <mwc-list-item @click=${this.handleRegisterBtn}>
        <mwc-button label="Register"></mwc-button>
      </mwc-list-item> `;
    let output: TemplateResult = login;
    let loginform = html``;

    if (this.showRegister) {
      output = html`<register-form
        @control-changed="${this.handleSnackbarMsg}"
      ></register-form>`;
    }

    if (this.showUserProfile) {
      output = html`
        <mwc-list-item>
          <mwc-icon slot="graphic">person</mwc-icon>
          ${this.username}</mwc-list-item
        >
        <mwc-list-item>${this.email}</mwc-list-item>
        <mwc-list-item @click=${this.handleLogout}>
          <mwc-button label="Logout"></mwc-button>
        </mwc-list-item>
      `;
    } else {
      loginform = html`<login-container> ${output} </login-container>`;
    }

    return html`
      <mwc-top-app-bar-fixed>
        <mwc-icon-button icon="menu" slot="navigationIcon"></mwc-icon-button>
        <div slot="title">Member Dashboard</div>
        <div slot="actionItems">${this.username}</div>
        <mwc-icon-button
          id="profileBtn"
          @click=${this.handleProfileClick}
          icon="person"
          slot="actionItems"
        ></mwc-icon-button>
        <mwc-menu id="menu" activatable> ${output} </mwc-menu>

        ${loginform}
        <member-stats></member-stats>

        <mwc-snackbar id="loginMessage" stacked labelText=${this.snackMessage}>
        </mwc-snackbar>
      </mwc-top-app-bar-fixed>
    `;
  }
}

customElements.define("top-bar", TopBar);
