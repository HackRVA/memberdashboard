import { LitElement, html, TemplateResult } from "lit-element";
import "./components/login-form";
import "./components/register-form";
import "@material/mwc-button";

console.log("this should load");
class MemberDashboard extends LitElement {
  showRegister: Boolean = false;
  handleRegisterBtn(): void {
    this.showRegister = !this.showRegister;
    this.requestUpdate()
  }

  render(): TemplateResult {
    if (this.showRegister) {
      return html`<register-form></register-form>`;
    }

    return html`
      <login-form></login-form>
      <mwc-button
        label="Register"
        @click=${this.handleRegisterBtn}
      ></mwc-button>
    `;
  }
}

const mountPoint: HTMLElement | null = document.getElementById(
  "memberdashboard"
);
customElements.define("member-dashboard", MemberDashboard);

/**
 * Mount the app into the DOM
 */
if (mountPoint)
  mountPoint.appendChild(document.createElement("member-dashboard"));
