// lit element
import { css, CSSResult } from "lit-element";

// membership
import { primaryDarkGreen } from "./colors";

export const registerFormStyles: CSSResult = css`
  mwc-formfield {
    display: block;
    margin-bottom: 16px;
  }

  mwc-button {
    width: 100%;
    --mdc-theme-primary: ${primaryDarkGreen};
  }
`;
