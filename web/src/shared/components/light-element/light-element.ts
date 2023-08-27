import { LitElement } from 'lit';

export class LightElement extends LitElement {
  createRenderRoot(): ShadowRoot | this {
    return this;
  }
}
