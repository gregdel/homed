const path = require("path");
const webpack = require("webpack");

const TerserPlugin = require("terser-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const WebpackPwaManifest = require("webpack-pwa-manifest");

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

module.exports = {
  mode: mode,
  entry: path.join(SRC_DIR, "js/app.js"),

  output: {
    publicPath: "",
    path: BUILD_DIR,
    filename: "[contenthash]-app.js",
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
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.svg$/,
        use: ["file-loader?name=[hash]-[name].[ext]"],
      },
    ],
  },

  optimization: {
    minimizer: [new TerserPlugin()],

    splitChunks: {
      cacheGroups: {
        vendors: {
          priority: -10,
          test: /[\\/]node_modules[\\/]/,
        },
      },

      chunks: "async",
      minChunks: 1,
      minSize: 30000,
      name: false,
    },
  },

  devtool: mode === "production" ? "source-map" : "inline-source-map",
};
