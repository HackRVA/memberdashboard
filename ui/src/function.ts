// jwt-decode
import jwt_decode from "jwt-decode";

// membership
import { JwtClaims } from "./components/shared/types";
import { Roles } from "./enums";

/**
 *
 * @param {string|Array} value
 * @description checks if the value is empty
 */
export const isEmpty = (value: string | any[]): boolean => {
  return value?.length === 0;
};

/**
 * @param {string} elementId - the #id of the component
 * @param {ShadowRoot} shadowRoot - the shadowRoot of the parent component
 * @description opens a component using its show() function
 */
export const showComponent = (
  elementId: string,
  shadowRoot: ShadowRoot
): void => {
  (
    shadowRoot?.querySelector(elementId) as HTMLElement & {
      show: Function;
    }
  ).show();
};

/**
 *
 * @description check if the user is an Admin
 */
export const isAdmin = (): boolean => {
  const jwt: string = localStorage.getItem("jwt");
  
  if (jwt) {
    const jwtClaims: JwtClaims = jwt_decode(jwt);
    
    const roles = new Set(jwtClaims.Groups);
    
    return roles.has(Roles.admin);
  }

  return false;
};
