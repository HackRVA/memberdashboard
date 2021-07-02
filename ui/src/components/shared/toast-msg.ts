// lit element
import {
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

// material
import { Snackbar } from "@material/mwc-snackbar";

// memberdashboard
import { ToastMessage } from "./types";

@customElement("toast-msg")
export class ToastMsg extends LitElement {
  @property({ type: Object })
  toastMsg: ToastMessage;

  toastMsgSnackBarTemplate: Snackbar;

  firstUpdated(): void {
    this.toastMsgSnackBarTemplate =
      this.shadowRoot?.querySelector("mwc-snackbar");
  }

  updated(): void {
    if (this.toastMsg) {
      this.toastMsgSnackBarTemplate.timeoutMs = this.toastMsg.duration;
    }
  }

  show(): void {
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
