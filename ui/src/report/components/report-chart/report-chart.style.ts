// lit element
import { css, CSSResult } from 'lit';

export const reportsChartStyle: CSSResult = css`
  #report-chart-container {
    margin-top: 56px;
    display: flex;
    align-items: center;
    flex-direction: column;
  }
  #membership-trends,
  #membership-distribution {
    margin-bottom: 36px;
    height: 560px;
    width: 700px;
  }
  .select-month {
    position: relative;
    float: right;
  }
  @media only screen and (max-width: 480px) {
    #membership-trends,
    #membership-distribution {
      height: 300px;
      width: 360px;
    }
  }
`;
