import { css, CSSResult } from "lit-element";

export const memberListStyles: CSSResult = css`
  h1 {
    margin-top: 0px;
    margin-bottom: 0px;
    justify-self: start;
  }
  .member-container {
    display: grid;
    justify-content: center;
    align-items: center;
    text-align: center;
    margin: 44px;
  }
  .member-header {
    display: inherit;
    grid-template-columns: 1fr 1fr;
    align-items: center;
  }
  .member-count {
    justify-self: end;
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
    border: 1px solid #e1e1e1;
    max-width: 320px;
  }
  table {
    margin-top: 24px;
    border-spacing: 0px;
  }
  .rfid-button {
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
    float: right;
    margin-top: 16px;
    width: 192px;
    margin-right: 24px;
    --mdc-typography-button-font-size: 12px;
  }

  .remove {
    --mdc-theme-primary: #e9437a;
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
`;
