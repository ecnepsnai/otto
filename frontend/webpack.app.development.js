const merge = require('webpack-merge');
const common = require('./webpack.app.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = merge.merge(common, {
    mode: "development",
    devtool: "source-map",
    plugins: [
        new HtmlWebpackPlugin({
            base: '/static/',
            template: './html/index.development.html',
            filename: 'index.html'
        }),
        new CopyPlugin({
            patterns: [
                { from: 'node_modules/react/umd/react.development.js', to: 'assets/js/' },
                { from: 'node_modules/react-dom/umd/react-dom.development.js', to: 'assets/js/' },
                { from: 'node_modules/react-router-dom/umd/react-router-dom.js', to: 'assets/js/' },
            ]
        }),
    ],
});
