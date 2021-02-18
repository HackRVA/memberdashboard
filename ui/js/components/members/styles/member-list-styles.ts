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
    grid-template-columns: 1fr 1fr 1fr;
    align-items: center;
  }
  .name {
    text-transform: capitalize;
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
  .remove {
    --mdc-theme-primary: #e9437a;
  }
  .horizontal-scrollbar {
    overflow: auto;
    max-width: 320px;
    white-space: nowrap;
  }
`;
