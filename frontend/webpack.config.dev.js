module.exports = {
  entry: './src/Index.tsx',
  output: {
    filename: '../backend/src/static/js/bundle.js'
  },

  devtool: '#cheap-module-eval-source-map',

  resolve: {
    extensions: ['.ts', '.tsx', '.js']
  },

  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: [
          {loader: 'ts-loader'}
        ]
      }
    ]
  }
};
