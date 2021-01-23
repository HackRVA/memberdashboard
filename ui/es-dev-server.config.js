const proxy = require('koa-proxies');

module.exports = {
  port: 8000,
  middlewares: [
    proxy('/edge', {
      target: 'http://localhost:3000',
      rewrite: path => path.replace(/\/edge/, '')
    }),
  ],
};