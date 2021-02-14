export const isEmpty = (value: string | any[]): boolean => {
  return value.length === 0;
};

export const showComponent = (
  elementId: string,
  shadowRoot: ShadowRoot
): void => {
  (shadowRoot?.querySelector(elementId) as HTMLElement & {
    show: Function;
  }).show();
};
