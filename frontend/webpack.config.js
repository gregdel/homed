const path = require("path");
const webpack = require("webpack");

const TerserPlugin = require("terser-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");

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
