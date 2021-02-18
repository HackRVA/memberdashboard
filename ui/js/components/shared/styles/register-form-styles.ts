import { css, CSSResult } from "lit-element";

export const registerFormStyles: CSSResult = css`
  .register-container {
    height: 270px;
    max-width: 250px;
    background-color: #e1e1e1;
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
