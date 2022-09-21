const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const ESLintPlugin = require('eslint-webpack-plugin');

let devtool = 'source-map';
let sourceType = 'development';
const version = process.env.VERSION;
const production = process.env.NODE_ENV === 'production';

if (production) {
    devtool = undefined;
    sourceType = 'production.min';

    if (!version) {
        throw new Error('VERSION environment variable must be specified for production builds');
    }
}

module.exports = {
    entry: './src/Login.tsx',

    mode: sourceType,
    devtool: devtool,

    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.html']
    },

    plugins: [
        new HtmlWebpackPlugin({
            base: production ? '/otto' + version + '/' : '/static/',
            template: './html/login.html',
            filename: 'login.html'
        }),
        new CopyPlugin({
            patterns: [
                { from: 'img/*.png', to: 'assets/', noErrorOnMissing: true },
                { from: 'img/*.svg', to: 'assets/', noErrorOnMissing: true },
                { from: 'img/*.jpg', to: 'assets/', noErrorOnMissing: true },
                { from: 'img/*.ico', to: 'assets/', noErrorOnMissing: true },
            ]
        }),
        new ESLintPlugin({
            extensions: ['.ts', '.tsx']
        }),
    ],

    module: {
        rules: [
            {
                test: /\.ts(x?)$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'ts-loader'
                    }
                ]
            },
            {
                test: /\.(woff|woff2)$/,
                type: 'asset/inline'
            },
            {
                enforce: 'pre',
                test: /\.js$/,
                loader: 'source-map-loader'
            },
            {
                test: /\.s[ac]ss$/i,
                use: [
                    {
                        loader: 'style-loader'
                    },
                    {
                        loader: 'css-loader'
                    },
                    {
                        loader: 'sass-loader'
                    }
                ]
            },
        ]
    },

    output: {
        path: path.resolve(__dirname, 'build'),
        hashFunction: 'xxhash64',
        filename: 'login.js'
    },
};
