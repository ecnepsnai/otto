const path = require('path');
const CopyPlugin = require('copy-webpack-plugin');

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
                    },
                    {
                        loader: "eslint-loader"
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
                test: /\.s[ac]ss$/i,
                use: [
                    // Creates `style` nodes from JS strings
                    'style-loader',
                    // Translates CSS into CommonJS
                    'css-loader',
                    // Compiles Sass to CSS
                    'sass-loader',
                ],
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
        'bootstrap': 'Bootstrap',
        'popper': 'Popper',
    },

    output: {
        path: path.resolve(__dirname, 'build'),
    },
};
