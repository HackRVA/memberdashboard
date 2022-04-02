// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import {
  primaryLightGray,
  primaryLightGreen,
  primaryRed,
} from '../../../shared/styles/colors';

export const resourceManagerStyle: CSSResult = css`
  .resource-container {
    display: grid;
    align-items: center;
    margin: 0px 16px;
    animation: fadeIn 1s;
  }
  .resource-header {
    display: inherit;
    grid-template-columns: 1fr 1fr;
    align-items: center;
  }
  .resource-header > h3 {
    justify-self: start;
  }
  .resource-header > div {
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
  .remove-acls {
    margin-left: 8px;
    margin-right: 32px;
    --mdc-theme-primary: ${primaryRed};
  }

  .online {
    color: ${primaryLightGreen};
  }
  .offline {
    color: ${primaryRed};
  }
`;
