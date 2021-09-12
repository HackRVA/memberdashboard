// lit element
import { css, CSSResult } from 'lit';
import { plainWhite, primaryDarkGray } from '../../styles';

// memberdashboard

export const memberDashboardContentStyle: CSSResult = css`
  .logout {
    margin-left: 24px;
    --mdc-theme-primary: ${plainWhite};
  }
  .version {
    float: right;
    color: ${primaryDarkGray};
  }
`;
