// lit element
import { css, CSSResult } from "lit-element";

// membership
import { primaryLightGray, primaryRed } from "./../../shared/styles/colors";

export const memberListStyles: CSSResult = css`
  .member-container {
    display: grid;
    align-items: center;
    text-align: center;
    margin: 20px;
  }
  .member-header {
    display: inherit;
    grid-template-columns: 1fr 1fr;
    align-items: center;
  }
  .member-count {
    justify-self: start;
  }
  .name {
    text-transform: capitalize;
  }
  .name span {
    position: relative;
    bottom: 13px;
  }
  td,
  th {
    text-align: left;
    padding: 8px;
    font-size: 20px;
    border: 1px solid ${primaryLightGray};
    max-width: 320px;
  }
  table {
    margin-top: 24px;
    border-spacing: 0px;
  }
  .buttons-container {
    justify-self: end;
  }

  .add-resource-to-members {
    float: right;
    margin-top: 16px;
    width: 225px;
    margin-right: 24px;
    --mdc-typography-button-font-size: 12px;
  }

  .refresh-members-list {
    margin-right: 24px;
  }
  
  .new-member-button {
    width: 192px;
    margin-right: 24px;
    --mdc-typography-button-font-size: 12px;
  }

  .remove {
    --mdc-theme-primary: ${primaryRed};
  }
  .horizontal-scrollbar {
    overflow: auto;
    max-width: 320px;
    white-space: nowrap;
  }

  .all-members-action-container {
    display: flex;
  }

  .all-members-checkbox {
    position: relative;
    top: 5px;
  }

  mwc-circular-progress {
    width: 20px;
    height: 20px;
    position: relative;
    bottom: 14px;
  }
`;
