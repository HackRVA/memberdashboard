// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import { plainWhite, primaryDarkGray } from '../../styles';

const appBarHeight: number = 64;
const tabHeight: number = 48;

export const mdContentStyle: CSSResult = css`
  mwc-tab-bar {
    background: var(--lumo-shade-color);
  }

  .logout {
    margin-left: 24px;
    --mdc-theme-primary: ${plainWhite};
  }
  .version {
    float: right;
    color: ${primaryDarkGray};
  }

  main {
    overflow: auto;
    position: absolute;
    top: ${appBarHeight + tabHeight}px;
    bottom: 0;
    left: 0;
    right: 0;
  }
`;
