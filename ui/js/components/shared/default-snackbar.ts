import { html, TemplateResult } from "lit-element";

export const defaultSnackbar = (
  elementId: string,
  message: string
): TemplateResult => {
  return html`
    <mwc-snackbar id=${elementId} stacked labelText=${message}> </mwc-snackbar>
  `;
};
