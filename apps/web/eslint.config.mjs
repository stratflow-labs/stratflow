// @ts-check

import js from "@eslint/js";
import tseslint from "typescript-eslint";
import reactHooks from "eslint-plugin-react-hooks";

export default tseslint.config(
  {
    ignores: [".next", "dist", "node_modules", "example", "src/shared/api/gen"],
  },

  js.configs.recommended,
  ...tseslint.configs.recommended,

  {
    files: ["**/*.{ts,tsx}"],
    plugins: {
      "react-hooks": reactHooks,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,

      "@typescript-eslint/no-unused-vars": [
        "warn",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
        },
      ],
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            {
              group: ["@/features/*/**"],
              message: "Import a feature through its public index.",
            },
            {
              group: ["@/widgets/*/**"],
              message: "Import a widget through its public index.",
            },
          ],
        },
      ],
    },
  },

  {
    files: ["src/shared/**/*.{ts,tsx}"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            {
              group: ["@/features/**", "@/widgets/**", "@/app/**"],
              message: "The shared layer cannot depend on upper layers.",
            },
          ],
        },
      ],
    },
  },
);
