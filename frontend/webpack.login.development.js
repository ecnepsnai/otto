const merge = require('webpack-merge');
const common = require('./webpack.login.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = merge.merge(common, {
    mode: 'development',
    devtool: 'source-map',
    plugins: [
        new HtmlWebpackPlugin({
            base: '/static/',
            template: './html/login.development.html',
            filename: 'login.html'
        }),
    ],
});
