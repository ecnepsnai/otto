const path = require('path');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
    entry: './src/Login.tsx',

    resolve: {
        extensions: [".ts", ".tsx", ".js", ".html"]
    },

    plugins: [
        new CopyPlugin({
            patterns: [
                { from: 'img/*.png', to: 'assets/', noErrorOnMissing: true},
                { from: 'img/*.svg', to: 'assets/', noErrorOnMissing: true},
                { from: 'img/*.jpg', to: 'assets/', noErrorOnMissing: true},
                { from: 'img/*.ico', to: 'assets/', noErrorOnMissing: true},
            ]
        })
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
            {
                enforce: "pre",
                test: /\.js$/,
                loader: "source-map-loader"
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
    externals: {
        'react': 'React',
        'react-dom': 'ReactDOM',
    },

    output: {
        path: path.resolve(__dirname, 'build'),
        filename: 'login.js'
    },
};
