// lit element
import { html, LitElement, TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';

// material
import { Snackbar } from '@material/mwc-snackbar';

// memberdashboard
import { ToastMessage } from '../../types/custom/toast-msg';
import { IPopup } from './../../types/custom/ipop-up';

@customElement('toast-msg')
export class ToastMsg extends LitElement implements IPopup {
  @property({ type: Object })
  toastMsg: ToastMessage;

  toastMsgSnackBarTemplate: Snackbar;

  firstUpdated(): void {
    this.toastMsgSnackBarTemplate =
      this.shadowRoot?.querySelector('mwc-snackbar');
  }

  updated(): void {
    if (this.toastMsg) {
      this.toastMsgSnackBarTemplate.timeoutMs = this.toastMsg.duration;
    }
  }

  public show(): void {
    this.toastMsgSnackBarTemplate.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-snackbar labelText=${this.toastMsg?.message}>
        <mwc-icon-button icon="close" slot="dismiss"></mwc-icon-button>
      </mwc-snackbar>
    `;
  }
}
