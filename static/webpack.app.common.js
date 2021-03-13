const path = require('path');
const CopyPlugin = require('copy-webpack-plugin');
const ESLintPlugin = require('eslint-webpack-plugin');

module.exports = {
    resolve: {
        extensions: [".ts", ".tsx", ".js", ".html"]
    },

    plugins: [
        new CopyPlugin({
            patterns: [
                { from: '404.html' },
                { from: 'img/*.png', to: 'assets/', noErrorOnMissing: true},
                { from: 'img/*.svg', to: 'assets/', noErrorOnMissing: true},
                { from: 'img/*.jpg', to: 'assets/', noErrorOnMissing: true},
                { from: 'img/*.ico', to: 'assets/', noErrorOnMissing: true},
            ]
        }),
        new ESLintPlugin({
            extensions: [".ts", ".tsx"]
        }),
    ],

    module: {
        rules: [
            {
                test: /\.ts(x?)$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: "ts-loader"
                    }
                ]
            },
            // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
            {
                enforce: "pre",
                test: /\.js$/,
                loader: "source-map-loader"
            },
            {
                test: /\.(woff|woff2)$/,
                use: {
                    loader: 'url-loader',
                },
            },
            {
                test: /\.s[ac]ss$/i,
                use: [{
                    // inject CSS to page
                    loader: 'style-loader'
                }, {
                    // translates CSS into CommonJS modules
                    loader: 'css-loader'
                }, {
                    // Run postcss actions
                    loader: 'postcss-loader'
                }, {
                    // compiles Sass to CSS
                    loader: 'sass-loader'
                }]
            },
        ]
    },

    // When importing a module whose path matches one of the following, just
    // assume a corresponding global variable exists and use that instead.
    // This is important because it allows us to avoid bundling all of our
    // dependencies, which allows browsers to cache those libraries between builds.
    externals: {
        'react': 'React',
        'react-dom': 'ReactDOM',
        'react-router-dom': 'ReactRouterDOM',
    },

    output: {
        path: path.resolve(__dirname, 'build'),
    },
};
