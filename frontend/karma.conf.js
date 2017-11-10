const args = process.argv;
args.splice(0, 4);

const polyfills = [];

var baseConfig = require('./karma.base.conf');

module.exports = (config) => {
  baseConfig(config);
  config.set({
    files: config.files.concat(args),
    browsers: ['ChromeHeadless']
  })
};
