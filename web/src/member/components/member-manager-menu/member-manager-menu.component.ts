// lit element
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// memberdashboard
import { memberManagerMenuStyle } from './member-manager-menu.style';
import dialogs from '../dialogs/dialogs';
import { coreStyle } from '../../../shared/styles';

@customElement('member-manager-menu')
export class MemberManagerMenu extends LitElement {
  @property({ type: Function })
  setActive;

  static get styles(): CSSResult[] {
    return [coreStyle, memberManagerMenuStyle];
  }

  render(): TemplateResult {
    return html`
      <div class="menuContainer">
        <vaadin-tabs>
          <vaadin-tab
            @click="${() => {
              this.setActive(true);
            }}"
          >
            Active
          </vaadin-tab>
          <vaadin-tab
            @click="${() => {
              this.setActive(false);
            }}"
          >
            Inactive
          </vaadin-tab>
        </vaadin-tabs>
        <vaadin-button @click=${dialogs(this.shadowRoot).editRFID(true)}>
          New Member
        </vaadin-button>
      </div>
    `;
  }
}
