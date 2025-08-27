import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";
import { lingui } from "@lingui/vite-plugin";
import importMetaUrlPlugin from "@codingame/esbuild-import-meta-url-plugin";
import vsixPluign from "@codingame/monaco-vscode-rollup-vsix-plugin";

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        react({
            babel: {
                plugins: ["@lingui/babel-plugin-lingui-macro"],
            },
        }),
        tailwindcss(),
        lingui(),
        vsixPluign(),
    ],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    optimizeDeps: {
        esbuildOptions: {
            plugins: [importMetaUrlPlugin],
        },
    },
    worker: {
        format: "es",
    },
});
