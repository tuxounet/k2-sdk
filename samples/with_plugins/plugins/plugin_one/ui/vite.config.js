import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";
// Resolve to the workspace root's node_modules to ensure single React instance
var rootNodeModules = path.resolve(__dirname, "../../../../node_modules");
var appFolder = path.basename(__dirname); // e.g. "horloge"
export default defineConfig({
    base: "./",
    plugins: [react()],
    build: {
        outDir: "dist",
        emptyOutDir: true,
        lib: {
            entry: path.resolve(__dirname, "src/Root.tsx"),
            name: "KH" + appFolder,
            formats: ["umd"],
            fileName: function () { return "app.umd.js"; },
        },
        rollupOptions: {
            external: [
                "react",
                "react-dom",
                "react/jsx-runtime",
                "react/jsx-dev-runtime",
                "react-router-dom",
                "@k-home/shared-ui",
            ],
            output: {
                // Inline all dynamic imports into a single chunk
                inlineDynamicImports: true,
                // Inject CSS into the JS bundle so consumers don't need a separate stylesheet
                assetFileNames: "app.[ext]",
                globals: {
                    react: "React",
                    "react-dom": "ReactDOM",
                    "react/jsx-runtime": "ReactJsxRuntime",
                    "react/jsx-dev-runtime": "ReactJsxDevRuntime",
                    "react-router-dom": "ReactRouterDOM",
                    "@k-home/shared-ui": "KHomeSharedUI",
                },
            },
        },
        // Produce a single self-contained file
        cssCodeSplit: false,
    },
    resolve: {
        alias: {
            react: path.resolve(rootNodeModules, "react"),
            "react-dom": path.resolve(rootNodeModules, "react-dom"),
            "react/jsx-runtime": path.resolve(rootNodeModules, "react/jsx-runtime"),
            "react/jsx-dev-runtime": path.resolve(rootNodeModules, "react/jsx-dev-runtime"),
        },
        dedupe: ["react", "react-dom", "react/jsx-runtime", "react-router-dom"],
        preserveSymlinks: false,
    },
    optimizeDeps: {
        include: ["react", "react-dom", "react-router-dom", "@k-home/shared-ui"],
        force: true,
    },
});
