// lit element
import { html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('abstract-toast')
export class AbstractToast extends LitElement {
  @property({ type: String })
  message: string;

  @property({ type: Number })
  timeoutMs: number = -1;

  firstUpdated(): void {
    this.shadowRoot?.querySelector('mwc-snackbar')?.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-snackbar
        labelText=${this.message}
        .timeoutMs=${this.timeoutMs}
        @MDCSnackbar:closed=${this.remove}
      >
        <mwc-icon-button icon="close" slot="dismiss"></mwc-icon-button>
      </mwc-snackbar>
    `;
  }
}

export function displayToast(message: string): void {
  const toast = document.createElement('abstract-toast');
  toast.message = message;
  toast.timeoutMs = 4000;
  document.body.appendChild(toast);
}
