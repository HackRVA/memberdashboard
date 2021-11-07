// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import { primaryBlue, primaryRed } from './colors';

export const coreStyle: CSSResult = css`
  .center-text {
    text-align: center;
  }

  .center {
    display: flex;
    justify-content: center;
  }

  .destructive-button {
    --mdc-theme-primary: ${primaryRed};
  }
  a {
    color: ${primaryBlue};
  }

  .margin-r-24 {
    margin-right: 24px;
  }

  /* ANIMATION */
  @keyframes fadeIn {
    0% {
      opacity: 0;
    }
    100% {
      opacity: 1;
    }
  }
`;
