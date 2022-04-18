// lit element
import { css, CSSResult } from 'lit';

// memberdashboard
import { primaryWhite } from '../../styles/colors';

export const mdCardStyle: CSSResult = css`
  /* On mouse-over, add a deeper shadow */
  .card:hover {
    box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.2);
  }
  /* Add some padding inside the card container */
  .container {
    padding: 24px;
  }
  .card {
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
    transition: 0.3s;
    border-radius: 5px; /* 5px rounded corners */
    min-height: 20vh;
    min-width: 50vw;
    overflow-x: scroll;
    background-color: var(--secondary-background-color);
  }
  card-container {
    display: grid;
    justify-content: center;
    margin-top: 5vh;
    margin-bottom: 5vh;
  }
  @media only screen and (max-width: 480px) {
    .container {
      padding: 32px 8px;
    }
  }
`;
