// lit element
import { customElement } from 'lit/decorators.js';
import { LitElement, TemplateResult, html, CSSResult } from 'lit';

// memberdashboard
import { coreStyle } from '../../styles/core.style';

@customElement('happy-minion')
export class HappyMinion extends LitElement {
  static get styles(): CSSResult[] {
    return [coreStyle];
  }

  // prettier-ignore
  render(): TemplateResult {
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
}
