const merge = require('webpack-merge');
const common = require('./webpack.login.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = merge.merge(common, {
    mode: "development",
    devtool: "source-map",
    plugins: [
        new HtmlWebpackPlugin({
            base: '/static/',
            template: './html/login.development.html',
            filename: 'login.html'
        }),
        new CopyPlugin({
            patterns: [
                { from: 'node_modules/react/umd/react.development.js', to: 'assets/js/' },
                { from: 'node_modules/react-dom/umd/react-dom.development.js', to: 'assets/js/' },
                { from: 'node_modules/bootstrap/dist/css/bootstrap.css', to: 'assets/css/' },
                { from: 'node_modules/bootstrap/dist/js/bootstrap.bundle.js', to: 'assets/js/' },
            ]
        }),
    ],
});
