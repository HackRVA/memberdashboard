// lit element
import { css, CSSResult } from 'lit';
import { primaryDarkGreen } from '../../../shared/styles/colors';

// memberdashboard

export const registerFormStyles: CSSResult = css`
  mwc-formfield {
    display: block;
    margin-bottom: 16px;
  }
  mwc-button {
    width: 100%;
    --mdc-theme-on-primary: var(--lumo-primary-color);
    --mdc-theme-primary: ${primaryDarkGreen};
  }

  :host {
    --mdc-theme-primary: var(--lumo-primary-color, rgb(0, 106, 245));
  }

  register-form {
    --mdc-theme-on-primary: var(--lumo-primary-color);
  }
`;
