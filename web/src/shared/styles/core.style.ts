// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import { primaryBlue, primaryRed, primaryGreen, primaryWhite } from './colors';

export const coreStyle: CSSResult = css`
  .text-center {
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

  .mr-24 {
    margin-right: 24px;
  }

  .mt-24 {
    margin-top: 24px;
  }

  .mr-8 {
    margin-right: 8px;
  }

  .ml-8 {
    margin-left: 8px;
  }

  .beta {
    color: ${primaryWhite};
    font-size: 14px;
    padding: 8px;
    line-height: 12px;
    background-color: ${primaryGreen};
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
