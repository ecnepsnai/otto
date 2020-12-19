const merge = require('webpack-merge');
const common = require('./webpack.app.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

const version = process.env.VERSION;
if (!version) {
    throw new Error('VERSION environment variable must be specified for production builds');
}

module.exports = merge.merge(common, {
    mode: "production",
    plugins: [
        new HtmlWebpackPlugin({
            base: '/otto' + version + '/',
            template: './html/index.production.html',
            templateParameters: {
                versionTag: version,
            },
            filename: 'index.html'
        }),
        new CopyPlugin({
            patterns: [
                { from: 'node_modules/react/umd/react.production.min.js', to: 'assets/js/' },
                { from: 'node_modules/react-dom/umd/react-dom.production.min.js', to: 'assets/js/' },
                { from: 'node_modules/react-router-dom/umd/react-router-dom.min.js', to: 'assets/js/' },
            ]
        }),
    ],
});
