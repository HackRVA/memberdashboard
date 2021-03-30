import { css, CSSResult } from "lit-element";

export const resourceManagerStyles: CSSResult = css`
  .resource-container {
    display: grid;
    align-items: center;
    margin: 44px;
  }

  .resource-header {
    display: inherit;
    grid-template-columns: 1fr 1fr;
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
    border: 1px solid #e1e1e1;
    max-width: 320px;
  }
  table {
    margin-top: 24px;
    border-spacing: 0px;
  }

  .remove {
    --mdc-theme-primary: #e9437a;
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
    background: #fdefbb;
    max-width: 200px;
    position: absolute;
    z-index: 1;
  }
`;
