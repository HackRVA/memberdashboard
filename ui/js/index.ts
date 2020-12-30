import { LitElement, html, TemplateResult } from "lit-element";
import { USER_PROFILE_ACTOR_ADDRESS } from "./constants";
import { UserActor } from "./actors/user";
import ActorStore from "./actors/actorStore";
import "./components/top-bar";

class MemberDashboard extends LitElement {
  firstUpdated(): void {
    ActorStore.register(USER_PROFILE_ACTOR_ADDRESS, UserActor);
  }
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
