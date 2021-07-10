// lit element
import { css, CSSResult } from "lit-element";

// memberdashboard
import { plainWhite, primaryDarkGray } from "./colors";

export const memberDashboardContentStyles: CSSResult = css`
  .logout {
    margin-left: 24px;
    --mdc-theme-primary: ${plainWhite};
  }

  .version {
    float: right;
    color: ${primaryDarkGray};
  }
`;
