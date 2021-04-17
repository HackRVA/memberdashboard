// lit element
import { css, CSSResult } from "lit-element";

// membership
import { primaryBlue, primaryWhite } from "./colors";

export const loginPageStyles: CSSResult = css`
  mwc-top-app-bar-fixed {
    --mdc-theme-primary: ${primaryWhite};
    --mdc-theme-on-primary: ${primaryBlue};
  }

  .login-container {
    display: grid;
    justify-content: center;
  }

  .text-center {
    text-align: center;
  }

  h1 {
    margin-top: 0px;
  }
`;
