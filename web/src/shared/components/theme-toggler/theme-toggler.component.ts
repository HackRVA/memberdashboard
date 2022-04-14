// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import { themeTogglerStyle } from './theme-toggler.style';

@customElement('theme-toggler')
export class ThemeToggler extends LitElement {
  @property()
  get darkModePreference(): string {
    const theme = JSON.parse(localStorage.getItem('dark-mode-preference'));

    if (theme == null) {
      return '';
    }

    return theme;
  }

  set darkModePreference(preference: string) {
    localStorage.setItem('dark-mode-preference', JSON.stringify(preference));
    document.documentElement.setAttribute('theme', preference);
  }

  static get styles(): CSSResult[] {
    return [themeTogglerStyle];
  }

  firstUpdated(): void {
    this.darkModePreference = JSON.parse(
      localStorage.getItem('dark-mode-preference')
    );
  }

  handleToggle = (): void => {
    if (this.darkModePreference !== 'dark') {
      this.darkModePreference = 'dark';
      return;
    }
    this.darkModePreference = '';
  };

  render(): TemplateResult {
    return html`
      <div class="theme-toggle-container">
        <vaadin-icon
          slot="prefix"
          @click=${this.handleToggle}
          icon="vaadin:lightbulb"
        ></vaadin-icon>
      </div>
    `;
  }
}
