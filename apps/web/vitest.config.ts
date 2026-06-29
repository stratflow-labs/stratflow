import { fileURLToPath } from "node:url";

import { defineConfig } from "vitest/config";

const root = fileURLToPath(new URL(".", import.meta.url));

export default defineConfig({
  root,
  resolve: {
    tsconfigPaths: true,
    alias: {
      "react-transition-group/TransitionGroupContext":
        "react-transition-group/cjs/TransitionGroupContext.js",
    },
  },
  test: {
    environment: "jsdom",
    setupFiles: ["src/test/setup.ts"],
    globals: true,
    css: true,
  },
});
