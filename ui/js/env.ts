export const ENV = {
  api: typeof process === "undefined" ? "/edge/api" : process.env.API,
};
