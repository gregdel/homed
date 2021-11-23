const path = require("path");
const webpack = require("webpack");

const HtmlWebpackPlugin = require("html-webpack-plugin");
const WebpackPwaManifest = require("webpack-pwa-manifest");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

var mode = "development";
if (process.env.NODE_ENV === "production") {
  mode = "production";
}

var SRC_DIR = path.resolve(__dirname, "src");
var BUILD_DIR = path.resolve(__dirname, "build");
if (process.env.NODE_ENV === "production") {
  mode = "production";
  BUILD_DIR = path.resolve(__dirname, "../build");
}

const config = {
  mode: mode,
  entry: path.join(SRC_DIR, "js/app.js"),

  output: {
    publicPath: "",
    path: BUILD_DIR,
    filename: "[contenthash]-app.js",
    assetModuleFilename: "[contenthash]-[name][ext][query]",
  },

  plugins: [
    new webpack.ProgressPlugin(),
    new HtmlWebpackPlugin({
      template: path.join(SRC_DIR, "html/index.html"),
    }),
    new WebpackPwaManifest({
      fingerprints: true,
      inject: true,
      ios: {
        "apple-mobile-web-app-status-bar-style": "default",
        "apple-mobile-web-app-title": "homed",
      },
      name: "homed",
      short_name: "homed",
      background_color: "#001529",
      theme_color: "#001529",
      display: "standalone",
      orientation: "omit",
      publicPath: "/",
      scope: "/",
      start_url: "/",
      icons: [
        {
          src: path.resolve(__dirname, "src/img/icon.svg"),
          sizes: [96, 128, 192, 256, 384, 512],
        },
        {
          src: path.resolve(__dirname, "src/img/icon.svg"),
          sizes: [80, 120, 152, 167, 180],
          ios: true,
        },
      ],
    }),
    new MiniCssExtractPlugin({
      filename: "[contenthash].css",
    }),
  ],

  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader",
      },
      {
        test: /\.css$/i,
        use: [MiniCssExtractPlugin.loader, "css-loader"],
      },
      {
        test: /\.svg$/,
        type: "asset",
      },
    ],
  },

  optimization: {},
  devtool: mode === "production" ? false : "source-map",
};

module.exports = config;
