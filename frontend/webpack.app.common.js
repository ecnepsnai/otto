const path = require('path');
const CopyPlugin = require('copy-webpack-plugin');
const ESLintPlugin = require('eslint-webpack-plugin');

module.exports = {
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.html']
    },

    plugins: [
        new CopyPlugin({
            patterns: [
                { from: '404.html' },
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
                        // inject CSS to page
                        loader: 'style-loader'
                    },
                    {
                        // translates CSS into CommonJS modules
                        loader: 'css-loader'
                    },
                    {
                        // compiles Sass to CSS
                        loader: 'sass-loader'
                    }
                ]
            },
        ]
    },

    output: {
        path: path.resolve(__dirname, 'build'),
        hashFunction: 'xxhash64',
    },
};
