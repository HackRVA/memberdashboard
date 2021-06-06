// lit element
import { css, CSSResult } from "lit-element";

// membership
import { primaryBlue, primaryDarkGreen } from "./colors";

export const loginFormStyles: CSSResult = css`
  mwc-formfield {
    display: block;
    margin-bottom: 24px;
  }

  mwc-button {
    margin-top: 32px;
    width: 100%;
    --mdc-theme-primary: ${primaryDarkGreen};
  }

  a {
    float: right;
    color: ${primaryBlue};
  }
`;
