const path = require('path');

module.exports = {
  mode: 'development',
  entry: {
    styles: './web/src/styles/finances.scss',
    bundle: './web/src/lib/component/main.tsx'
  },
  output: {
    path: path.resolve(__dirname, 'web', 'dist'),
    filename: '[name].js',
    publicPath: '/finances/scripts/bundle.js'
  },
  devtool: 'source-map',
  module: {
    rules: [
      {
        test: /\.scss$/,
        use: [
          {loader: 'file-loader', options: {name: 'finances.css'}},
          {loader: 'extract-loader'},
          {loader: 'css-loader'},
          {loader: 'sass-loader', options: {sassOptions: {includePaths: ['./node_modules']}}},
        ]
      },
      {
        test: /\.tsx?$/,
        loader: 'babel-loader',
        include: /web\/src\/lib/,
      }
    ]
  },
  resolve: {
    modules: [path.resolve(__dirname, 'public'), 'node_modules'],
    extensions: ['.js', '.jsx', '.ts', '.tsx'],
  },
};