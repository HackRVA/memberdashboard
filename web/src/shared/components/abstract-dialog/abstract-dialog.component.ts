// lit element
import { html, CSSResult, LitElement, TemplateResult, render } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { guard } from 'lit/directives/guard.js';

// memberdashboard
import { IPopup } from '../../types/custom/ipop-up';
import { abstractDialogStyle } from './abstract-dialog.style';
import { coreStyle } from '../../styles';

@customElement('abstract-dialog')
export class AbstractDialog extends LitElement implements IPopup {
  @property()
  heading: string;
  @property()
  dialogLayout: TemplateResult;

  @state()
  dialogOpened: boolean = true;

  show(): void {
    this.dialogOpened = true;
  }

  hide(): void {
    this.dialogOpened = false;
  }

  render(): TemplateResult {
    return html`<vaadin-dialog
      .opened="${this.dialogOpened}"
      @opened-changed="${(e: CustomEvent) => {
        this.dialogOpened = e.detail.value

        if (!this.dialogOpened) this?.remove(); /* clean up the DOM */
      }}"
      .renderer="${guard([], () => (root: HTMLElement) => {
        render(html`
        <div class="dialog-container">
          <h1> ${this.heading} </h1>
          ${this.dialogLayout}
        </div>
      `, root);
      })}"
    ></vaadin-dialog>
    `;
  }
}
