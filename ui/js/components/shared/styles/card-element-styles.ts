import { css, CSSResult } from "lit-element";

export const cardElementStyles: CSSResult = css`
  /* On mouse-over, add a deeper shadow */
  .card:hover {
    box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.2);
  }

  /* Add some padding inside the card container */
  .container {
    padding: 2px 16px;
  }

  .card {
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
    transition: 0.3s;
    border-radius: 5px; /* 5px rounded corners */
    min-height: 20vh;
    min-width: 60vw;
    background-color: whitesmoke;
  }
  card-container {
    width: 100vw;
    display: grid;
    justify-content: center;
    margin-top: 5vh;
    margin-bottom: 5vh;
  }
`;
