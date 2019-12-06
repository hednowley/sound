const common = require("./webpack.common");
const merge = require("webpack-merge");
const path = require("path");

module.exports = merge(common(false), {
  mode: "development",
  devtool: "inline-source-map",
  devServer: {
    contentBase: path.join(__dirname, "dist"),
    port: 9000,
    historyApiFallback: true,
    proxy: {
      "/api": {
        target: "http://localhost:3684"
      },
      "/ws": {
        target: "http://localhost:3684",
        ws: true
      }
    }
  }
});
