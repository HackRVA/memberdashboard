import {
  customElement,
  html,
  LitElement,
  property,
  TemplateResult,
} from "lit-element";

@customElement("loading-content")
export class LoadingContent extends LitElement {
  @property({ type: Boolean })
  finishedLoading: boolean;

  displayPage(): TemplateResult {
    if (this.finishedLoading) {
      return html` <slot> </slot> `;
    }

    return html`
      <mwc-circular-progress indeterminate></mwc-circular-progress>
    `;
  }

  render(): TemplateResult {
    return html` ${this.displayPage()} `;
  }
}
