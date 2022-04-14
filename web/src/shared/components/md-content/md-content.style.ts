// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import { plainWhite, primaryDarkGray } from '../../styles';

export const mdContentStyle: CSSResult = css`
  .logout {
    margin-left: 24px;
    --mdc-theme-primary: ${plainWhite};
  }
  .version {
    float: right;
    color: ${primaryDarkGray};
  }

  .theme-toggle-container {
    display: flex;
    padding-right: 2rem;
  }
`;
