// lit element
import { customElement, property, state } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';

// vaadin
import { TextFieldValueChangedEvent } from '@vaadin/text-field';

@customElement('member-search')
export class MemberSearch extends LitElement {
  @property({ type: Function })
  memberSearch;

  private input: string;

  constructor() {
    super();
    this.addEventListener('keypress', (e: KeyboardEvent) => {
      if (e.key !== 'Enter') return;
      this.search();
    });
  }

  search = (): void => {
    const searchTerm = ((this.input as string) || '').trim();
    this.memberSearch(searchTerm);
  };

  render(): TemplateResult {
    return html`
      <vaadin-text-field
        placeholder="Search"
        style="width: 50%;"
        @value-changed="${(e: TextFieldValueChangedEvent) =>
          (this.input = e.detail.value)}"
      >
        <vaadin-icon slot="prefix" icon="vaadin:search"></vaadin-icon>
      </vaadin-text-field>
      <vaadin-button @click="${this.search}">search</vaadin-button>
    `;
  }
}
