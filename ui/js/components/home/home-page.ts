// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
} from "lit-element";

// membership
import { homePageStyles } from "./styles/home-page-styles";
import "../shared/card-element";

@customElement("home-page")
export class HomePage extends LitElement {
  static get styles(): CSSResult[] {
    return [homePageStyles];
  }

  firstUpdated(): void {}
  // prettier-ignore
  displayHomePage(): TemplateResult {
    return html`
    <div class="center">
      <pre>
      ──────────▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
      ────────█═════════════════█
      ──────█═════════════════════█
      ─────█═══▄▄▄▄▄▄▄═══▄▄▄▄▄▄▄═══█
      ────█═══█████████═█████████═══█
      ────█══██▀────▀█████▀────▀██══█
      ───██████──█▀█──███──█▀█──██████
      ───██████──▀▀▀──███──▀▀▀──██████
      ────█══▀█▄────▄██─██▄────▄█▀══█
      ────█════▀█████▀───▀█████▀════█
      ────█═════════════════════════█
      ────█══════▀▄▄▄▄▄▄▄▄▄▄▄═══════█
      ───▐▓▓▌═════════════════════▐▓▓▌
      ───▐▐▓▓▌▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▐▓▓▌▌
      ───█══▐▓▄▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄▓▌══█
      ──█══▌═▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▌═▐══█
      ──█══█═▐▓▓▓▓▓▓▄▄▄▄▄▄▄▓▓▓▓▓▓▌═█══█
      ──█══█═▐▓▓▓▓▓▓▐██▀██▌▓▓▓▓▓▓▌═█══█
      ──█══█═▐▓▓▓▓▓▓▓▀▀▀▀▀▓▓▓▓▓▓▓▌═█══█
      ──█══█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█══█
      ─▄█══█▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▌█══█▄
      ─█████▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▌─█████
      ─██████▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▌─██████
      ──▀█▀█──▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▌───█▀█▀
      ─────────▐▓▓▓▓▓▓▌▐▓▓▓▓▓▓▌
      ──────────▐▓▓▓▓▌──▐▓▓▓▓▌
      ─────────▄████▀────▀████▄
      ─────────▀▀▀▀────────▀▀▀▀
      </pre>
    <div>
    `;
  }

  render(): TemplateResult {
    return html` <card-element> ${this.displayHomePage()} </card-element> `;
  }
}
