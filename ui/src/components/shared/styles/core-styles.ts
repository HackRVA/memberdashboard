import { css, CSSResult } from "lit-element";
import { primaryBlue, primaryRed } from "./colors";

export const coreStyles: CSSResult = css`
  .center-text {
    text-align: center;
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
`;
