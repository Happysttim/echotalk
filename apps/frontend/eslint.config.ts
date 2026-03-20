import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import { defineConfig } from "eslint/config";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default defineConfig([
  {
    files: ["**/*.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        projectService: {
          allowDefaultProject: [
            "*.ts",
            "*.mts",
            "eslint.config.ts",
            "eslint.config.mts",
          ],
        },
        tsconfigRootDir: __dirname,
      },
    },
  },

  tseslint.configs.recommended,
  tseslint.configs.recommendedTypeChecked,

  {
    ...pluginReact.configs.flat.recommended,
    files: ["**/*.{jsx,tsx}"],
    languageOptions: {
      ...pluginReact.configs.flat.recommended?.languageOptions,
      globals: {
        ...globals.browser,
      },
    },
  },
]);