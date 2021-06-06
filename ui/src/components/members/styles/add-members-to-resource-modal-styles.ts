import { css, CSSResult } from "lit-element";

export const addMembersToResourceModalStyles: CSSResult = css`
  .emails {
    text-align: center;
    margin-bottom: 16px;
    max-height: 300px;
    overflow-y: scroll;
  }

  mwc-dialog {
    --mdc-dialog-min-width: 400px;
  }

  mwc-select {
    width: 400px;
  }
`;
