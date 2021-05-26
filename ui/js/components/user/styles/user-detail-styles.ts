import { css, CSSResult } from "lit-element";

export const userDetailStyles: CSSResult = css`
  .user-profile {
    display: flex;
    gap: calc(100% / 8);
    justify-content: center;
  }

  paper-card {
    height: 320px;
    width: 360px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
  }

  dl {
    font-size: 16px;
    text-align: left;
  }

  dt {
    margin-bottom: 8px;
    font-weight: bold;
  }

  dd {
    margin-bottom: 8px;
  }

  li {
    margin-bottom: 8px;
  }

  .card-actions {
    text-align: center;
  }

  .lenny-face {
    text-align: center;
    font-size: 56px;
    margin-top: 32px;
  }

  @media only screen and (max-width: 480px) {
    .user-profile {
      display: flex;
      flex-direction: column;
    }

    paper-card {
      margin-bottom: 24px;
    }
  }
`;
