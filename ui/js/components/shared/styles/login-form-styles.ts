// lit element
import { css, CSSResult } from "lit-element";

// membership
import { primaryGray } from "./colors";

export const loginFormStyles: CSSResult = css`
  .login-container {
    height: 270px;
    max-width: 250px;
    background-color: ${primaryGray};
    padding: 24px;
  }
  mwc-formfield {
    display: block;
    margin-bottom: 16px;
  }
  .mwc-button {
    margin-bottom: 12px;
  }
  .register {
    float: left;
  }
  .login {
    float: right;
  }
`;
