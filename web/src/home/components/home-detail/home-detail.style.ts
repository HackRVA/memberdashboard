// lit element
import { css, CSSResult } from 'lit';

export const homeDetailStyle: CSSResult = css`
  #home-detail {
    display: flex;
    align-items: center;
    text-align: center;
    flex-direction: column;
  }

  mwc-circular-progress {
    height: 24px;
    width: 24px;
    position: relative;
    bottom: 12px;
    left: 4px;
  }
`;
