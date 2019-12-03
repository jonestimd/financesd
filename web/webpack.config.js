const path = require('path');
const autoprefixer = require('autoprefixer');

module.exports = [{
    entry: {
      styles: './src/styles/app.scss',
      bundle: './src/lib/main.tsx'
    },
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: '[name].js',
      publicPath: '/finances/scripts/bundle.js'
    },
    devtool: 'source-map',
    module: {
      rules: [
        {
          test: /\.scss$/,
          use: [
            { loader: 'file-loader', options: {name: 'finances.css'} },
            { loader: 'extract-loader' },
            { loader: 'css-loader' },
            // { loader: 'postcss-loader', options: {plugins: () => [autoprefixer()]} },
            { loader: 'sass-loader' },
          ]
        },
        {
          test: /\.tsx?$/,
          loader: 'babel-loader',
          include: /src\/lib/,
        }
      ]
    },
    resolve: {
      modules: [path.resolve(__dirname, 'public'), 'node_modules'],
      extensions: ['.js', '.jsx', '.ts', '.tsx'],
    },
  }];