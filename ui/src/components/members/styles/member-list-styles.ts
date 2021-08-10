// lit element
import { css, CSSResult } from "lit-element";

// memberdashboard
import {
  primaryBlue,
  primaryLightGray,
  primaryRed,
} from "../../shared/styles/colors";

export const memberListStyles: CSSResult = css`
  .member-container {
    display: grid;
    align-items: center;
    text-align: center;
    margin: 0px 16px;
    animation: fadeIn 1s;
  }

  .member-header {
    display: inherit;
    grid-template-columns: 1fr 1fr;
    align-items: center;
  }

  .member-header > h3 {
    justify-self: start;
  }

  .member-header > div {
    justify-self: end;
  }

  table {
    margin-top: 24px;
    border-spacing: 0px;
  }

  td,
  th {
    text-align: left;
    padding: 8px;
    font-size: 20px;
    border: 1px solid ${primaryLightGray};
    max-width: 320px;
  }

  td:first-child {
    text-transform: capitalize;
  }

  td:first-child > span {
    position: relative;
    bottom: 13px;
  }

  td:last-child > mwc-icon-button {
    color: ${primaryBlue};
  }

  .all-members-action-container {
    display: flex;
  }

  .all-members-action-container > mwc-button {
    float: right;
    margin-top: 16px;
    width: 225px;
    margin-right: 24px;
    --mdc-typography-button-font-size: 12px;
  }

  .all-members-action-container > mwc-formfield {
    position: relative;
    top: 5px;
  }

  .more-actions-container {
    position: relative;
  }

  .add-resources {
    color: ${primaryBlue};
  }

  .remove-resources {
    color: ${primaryRed};
  }

  .horizontal-scrollbar {
    overflow: auto;
    max-width: 320px;
    white-space: nowrap;
  }
`;
