// lit element
import { css, CSSResult } from "lit-element";

// membership
import {
  primaryGray,
  primaryGreen,
  primaryRed,
  primaryYellow,
} from "./../../shared/styles/colors";

export const resourceManagerStyles: CSSResult = css`
  .resource-container {
    display: grid;
    align-items: center;
    margin: 44px;
  }

  .resource-header {
    display: inherit;
    grid-template-columns: 1fr;
    align-items: center;
  }

  .button-container {
    justify-self: end;
  }

  td,
  th {
    text-align: left;
    padding: 8px;
    font-size: 20px;
    border: 1px solid ${primaryGray};
    max-width: 320px;
  }
  table {
    margin-top: 24px;
    border-spacing: 0px;
  }

  .remove {
    --mdc-theme-primary: ${primaryRed};
  }

  .mr-8 {
    margin-right: 8px;
  }

  .mr-32 {
    margin-right: 32px;
  }

  .update-acls .note {
    visibility: hidden;
  }

  .update-acls:hover .note {
    visibility: visible;
  }

  .note {
    margin-top: 8px;
    padding: 8px;
    font-size: 14px;
    background: ${primaryYellow};
    max-width: 200px;
    position: absolute;
    z-index: 1;
  }

  .online {
    color: ${primaryGreen};
  }

  .offline {
    color: ${primaryRed};
  }
`;
