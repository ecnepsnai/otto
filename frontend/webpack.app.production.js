const merge = require('webpack-merge');
const common = require('./webpack.app.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const version = process.env.VERSION;
if (!version) {
    throw new Error('VERSION environment variable must be specified for production builds');
}

module.exports = merge.merge(common, {
    mode: 'production',
    plugins: [
        new HtmlWebpackPlugin({
            base: '/otto' + version + '/',
            template: './html/index.production.html',
            templateParameters: {
                versionTag: version,
            },
            filename: 'index.html'
        }),
    ],
});
