// lit element
import { customElement } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import '../../../shared/components/md-card';
import { homePageStyle } from './home-page.style';

@customElement('home-page')
export class HomePage extends LitElement {
  static get styles(): CSSResult[] {
    return [homePageStyle];
  }
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
    return html` <md-card> ${this.displayHomePage()} </md-card> `;
  }
}
