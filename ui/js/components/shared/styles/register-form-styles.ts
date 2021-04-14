// lit element
import { css, CSSResult } from "lit-element";

// membership
import { primaryGray } from "./colors";

export const registerFormStyles: CSSResult = css`
  .register-container {
    height: 270px;
    max-width: 250px;
    background-color: ${primaryGray};
    padding: 24px;
  }
  mwc-formfield {
    display: block;
    margin-bottom: 16px;
  }
  .sign-in {
    float: left;
  }
  mwc-button {
    float: right;
    margin-bottom: 12px;
  }
`;
