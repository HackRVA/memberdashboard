// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import {
  primaryBlue,
  primaryWhite,
  plainWhite,
  primaryDarkGray,
} from '../../../shared/styles/colors';

export const loginPageStyle: CSSResult = css`
  .login-container {
    display: grid;
    justify-content: center;
    align-content: center;
    height: 400px;
  }
  login-form,
  register-form {
    height: 250px;
    padding: 24px 16px;
    background-color: var(--secondary-background-color)};
    border-radius: 8px;
    border: 1px solid silver;
  }

  .toggle-form-text {
    margin-top: 16px;
    padding: 12px;
    border: 1px solid ${primaryWhite};
    border-radius: 8px;
    opacity: 0.7;
  }

  :host {
    --mdc-theme-primary: var(--lumo-primary-color, rgb(0, 106, 245));
  }

  login-form {
    --mdc-theme-on-primary: var(--lumo-primary-color);
  }
`;
