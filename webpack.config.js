const VueLoaderPlugin = require('vue-loader/lib/plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')
const CleanWebpackPlugin = require('clean-webpack-plugin')

module.exports = (_, argv) => {
  let mode = 'development';
  if (argv !== undefined && 'mode' in argv) {
    mode = argv.mode;
  }
  const prod = mode === 'production';

  return {
    mode: mode,
    devtool: prod ? false : 'source-map',
    stats: {children: false},
    watchOptions: {
      ignored: ['node_modules', 'static', '__tests__', 'app'],
    },
    context: __dirname + '/resources/assets',
    entry: {
      'fingerprints/build/fingerprints': './fingerprints/js/main.js',
      'locations/build/locations': './locations/js/main.js',
      'corpmappings/build/corpmappings': './corpmappings/js/main.js',
      'tags/build/tags': './tags/js/main.js',
      'customerusers/build/customerusers': './customerusers/customerusers.css',
    },
    output: {
      path: __dirname + '/static',
      filename: prod ? '[name]-[contentHash].js' : '[name].js',
    },
    module: {
      rules: [
        {
          test: /\.vue$/,
          loader: 'vue-loader',
        },
        {
          test: /\.js$/,
          exclude: /(node_modules)/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['@babel/preset-env'],
              plugins: ['@babel/plugin-transform-runtime']
            }
          }
        },
        {
          test: /\.scss$/,
          use: [
            {loader: MiniCssExtractPlugin.loader},
            {loader: 'css-loader', options: {sourceMap: !prod}},
            {loader: 'postcss-loader', options: {config: {ctx: {mode}}}},
            {loader: 'sass-loader', options: {sourceMap: !prod}},
          ],
        },
        {
          test: /\.(png|woff|woff2|eot|ttf|svg)$/,
          loader: 'file-loader',
          options: {
            publicPath: '/static/',
            outputPath: 'public/',
            name: 'build/fonts/[name]-[contentHash].[ext]'
          }
        },
        {
          test: /\.css$/,
          use: [
            MiniCssExtractPlugin.loader,
            'css-loader',
            'postcss-loader'
          ]
        }
      ],
    },
    resolve: {
      extensions: ['.js', '.vue'],
    },
    plugins: [
      new CleanWebpackPlugin(['static/fingerprints/build']),
      new CleanWebpackPlugin(['static/locations/build']),
      new CleanWebpackPlugin(['static/corpmappings/build']),
      new CleanWebpackPlugin(['static/tags/build']),
      new CleanWebpackPlugin(['static/customerusers/build']),
      new VueLoaderPlugin(),
      new MiniCssExtractPlugin({
        filename: prod ? '[name]-[contentHash].css' : '[name].css',
        chunkFilename: '[id].css',
      }),
      new ManifestPlugin({
        fileName: __dirname + '/static/manifest.json',
        serialize(manifest) {
          const stripped = {}
          Object.keys(manifest).forEach(key => {
            const k = key.replace(/^(public|fingerprints|locations|corpmappings|tags|customerusers)\/build\//, '')
            stripped[k] = manifest[key].replace(/^(public|fingerprints|locations|corpmappings|tags|customerusers)\//, '')
          })
          return JSON.stringify(stripped, null, 2)
        }
      })
    ],
  }
}
