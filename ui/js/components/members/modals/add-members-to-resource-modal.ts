import { Dialog } from "@material/mwc-dialog";
import "@material/mwc-dialog";
import {
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

@customElement("add-members-to-resource-modal")
export class AddMembersToResourceModal extends LitElement {
  @property({ type: Array })
  emails: string[] = [];

  addResourceToMembersModalTemplate: Dialog;

  firstUpdated(): void {
    this.addResourceToMembersModalTemplate = this.shadowRoot.querySelector(
      "mwc-dialog"
    );
  }

  show(): void {
    this.addResourceToMembersModalTemplate?.show();
  }

  render(): TemplateResult {
    return html`
      <mwc-dialog heading="Assign resource to members">
        <div>
          ${this.emails?.map((x: string) => {
            return html`<div>${x}</div>`;
          })}
        </div>
        <mwc-select label="Resources">
          <mwc-list-item> x </mwc-list-item>
        </mwc-select>
      </mwc-dialog>
    `;
  }
}
