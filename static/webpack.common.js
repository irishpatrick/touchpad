const path = require("path")
const ESLintPlugin = require("eslint-webpack-plugin")

module.exports = {
    entry: "./src/index.js",
    plugins: [new ESLintPlugin({})],
    output: {
        filename: "main.bundle.js",
        path: path.resolve(__dirname, "dist"),
        clean: false,
    },
}

