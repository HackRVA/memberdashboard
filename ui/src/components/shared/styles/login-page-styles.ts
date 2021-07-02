// lit element
import { css, CSSResult } from "lit-element";

// memberdashboard
import {
  primaryBlue,
  primaryWhite,
  plainWhite,
  primaryDarkGray,
} from "./colors";

export const loginPageStyles: CSSResult = css`
  mwc-top-app-bar-fixed {
    --mdc-theme-primary: ${primaryWhite};
    --mdc-theme-on-primary: ${primaryBlue};
  }

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
    background-color: ${primaryDarkGray};
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
    background-color: ${primaryDarkGray};
    opacity: 0.7;
    color: ${plainWhite};
  }
`;
