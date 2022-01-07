// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import { primaryDarkGreen } from '../../../shared/styles';

export const loginFormStyle: CSSResult = css`
  mwc-formfield {
    display: block;
    margin-bottom: 24px;
  }
  mwc-button {
    margin-top: 32px;
    width: 100%;
    --mdc-theme-primary: ${primaryDarkGreen};
  }
`;
