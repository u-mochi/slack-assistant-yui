const polyfills = [];

module.exports = (config) => {
  config.set({
    basePath: '',
    frameworks: ['jasmine'],
    files: polyfills,
    preprocessors: {
      '**/*.spec.ts': ['webpack'],
      '**/*.spec.tsx': ['webpack']
    },
    mime: {
      'text/x-typescript': ['ts','tsx']
    },
    webpack: {
      resolve: {
        extensions: ['.ts', '.js', ".tsx"]
      },
      module: {
        rules: [
          {
            test: /\.tsx?$/,
            use: [
              {loader: "ts-loader"}
            ]
          }
        ]
      },
      externals: {
        'react/lib/ExecutionEnvironment': true,
        'react/addons': true,
        'react-addons-test-utils': true,
        'react/lib/ReactContext': 'window',
        'fs': '{}'
      }
    },
    webpackMiddleware: {
      stats: 'errors-only',
      noInfo: true
    },
    reporters: ['mocha'],
    port: 9876,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: false,
    singleRun: true,
    concurrency: Infinity
  })
};
