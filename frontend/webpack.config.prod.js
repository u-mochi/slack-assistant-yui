const webpack = require('webpack');

module.exports = {
  entry: "./src/Index.tsx",
  output: {
    filename: "../backend/src/static/js/bundle.js"
  },

  devtool: '#source-map',

  resolve: {
    extensions: [".ts", ".tsx", ".js"]
  },

  plugins: [
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      },
      sourceMap: true
    }),

    new webpack.DefinePlugin({
      'process.env':{
        'NODE_ENV': JSON.stringify('production')
      }
    })
  ],
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: [
          {loader: "ts-loader"}
        ]
      }
    ]
  }
};
