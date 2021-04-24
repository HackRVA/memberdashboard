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
    align-content: center;
    height: 40vh;
  }

  login-form,
  register-form {
    height: 250px;
    padding: 24px;
    background-color: #888888;
    border-radius: 8px;
    border: 1px solid silver;
  }

  .text-center {
    text-align: center;
  }

  .toggle-form-text {
    margin-top: 16px;
    padding: 12px;
    border: 1px solid ${primaryWhite};
    border-radius: 8px;
    background-color: #888888;
    opacity: 0.7;
    color: white;
  }

  a {
    color: ${primaryBlue};
  }
`;
