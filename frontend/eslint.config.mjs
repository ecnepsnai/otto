import typescriptEslint from "@typescript-eslint/eslint-plugin";
import tsParser from "@typescript-eslint/parser";
import path from "node:path";
import { fileURLToPath } from "node:url";
import js from "@eslint/js";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

export default [{
    ignores: [
        "**/node_modules/",
        "**/build/",
        "**/webpack.*.js",
        "**/start_webpack.js",
    ],
}, ...compat.extends(
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react/recommended",
), {
    plugins: {
        "@typescript-eslint": typescriptEslint,
    },

    languageOptions: {
        parser: tsParser,
    },

    settings: {
        react: {
            createClass: "createReactClass",
            pragma: "React",
            version: "detect",
            flowVersion: "0.53",
        },

        propWrapperFunctions: ["forbidExtraProps", {
            property: "freeze",
            object: "Object",
        }, {
            property: "myFavoriteWrapper",
        }],

        linkComponents: ["Hyperlink", {
            name: "Link",
            linkAttribute: "to",
        }],
    },
}, {
    files: ["**/*.js", "**/*.ts", "**/*.tsx"],

    rules: {
        quotes: ["error", "single"],
        "brace-style": ["error", "1tbs"],
        "no-var": "error",
        "react/prop-types": "off"
    },
}, {
    files: ["**/*.ts", "**/*.tsx"],

    rules: {
        "@typescript-eslint/no-namespace": "off",
        semi: "error",
    },
}, {
    files: ["src/types/*.ts"],

    rules: {
        "@typescript-eslint/explicit-module-boundary-types": "off",
        "@typescript-eslint/no-explicit-any": "off",
    },
}];