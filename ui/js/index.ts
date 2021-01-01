import { LitElement, html, TemplateResult } from "lit-element";
import "./components/top-bar";

class MemberDashboard extends LitElement {
  render(): TemplateResult {
    return html` <top-bar></top-bar> `;
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
