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
  (shadowRoot?.querySelector(elementId) as HTMLElement & {
    show: Function;
  }).show();
};
