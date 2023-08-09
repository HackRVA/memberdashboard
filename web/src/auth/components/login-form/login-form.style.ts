// lit element
import { css, CSSResult } from 'lit';

export const loginFormStyle: CSSResult = css`
  @import url(https://fonts.googleapis.com/css?family=Roboto);
  html {
    font-family: 'Roboto', sans-serif;
  }

  form {
    display: grid;
    justify-content: center;
  }

  .form__group {
    position: relative;
    padding: 15px 0 0;
    margin-top: 10px;
  }

  .form__field {
    font-family: inherit;
    width: 100%;
    border: 0;
    border-bottom: 1px solid #d2d2d2;
    outline: 0;
    font-size: 16px;
    color: rgb(0, 95, 219);
    padding: 7px 0;
    background: transparent;
    transition: border-color 0.2s;
  }

  .form__field::placeholder {
    color: transparent;
  }

  .form__field:placeholder-shown ~ .form__label {
    font-size: 16px;
    cursor: text;
    top: 20px;
  }

  label,
  .form__field:focus ~ .form__label {
    position: absolute;
    top: 0;
    display: block;
    transition: 0.2s;
    font-size: 12px;
    color: rgb(0, 95, 219);
  }

  .form__field:focus ~ .form__label {
    color: rgb(0, 95, 219);
  }

  .form__field:focus {
    padding-bottom: 6px;
    border-bottom: 2px solid rgb(0, 95, 219);
  }
`;
