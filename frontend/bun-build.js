#!/usr/bin/env bun

// NOTE: This must be run with: NODE_ENV=production bun run bun-build.js
// Or use the package.json script which sets it automatically

import { copyFileSync, mkdirSync, rmSync, existsSync } from "fs";
import { join } from "path";

const args = process.argv.slice(2);
const isDev = args.includes("--dev");
const outdir = isDev ? "build" : "../build";

// Clean and create output directory
if (existsSync(outdir)) {
  rmSync(outdir, { recursive: true });
}
mkdirSync(outdir, { recursive: true });

// Copy static assets
const assets = [
  "app.css",
  "icon.png",
  "icon.svg",
  "index.html",
  "manifest.json",
  "service-worker.js",
];

for (const asset of assets) {
  const srcPath = join("src/assets", asset);
  const destPath = join(outdir, asset);
  if (existsSync(srcPath)) {
    copyFileSync(srcPath, destPath);
  }
}

// Build JavaScript/TypeScript with Bun's bundler
// Bun automatically handles .ts/.tsx files
const result = await Bun.build({
  entrypoints: ["./src/js/app.tsx"],
  outdir: outdir,
  minify: !isDev,
  sourcemap: isDev ? "external" : "none",
  target: "browser",
  format: "iife",
  splitting: false,
  naming: "[dir]/[name].[ext]",
});

if (!result.success) {
  console.error("Build failed:");
  for (const message of result.logs) {
    console.error(message);
  }
  process.exit(1);
}

const totalSize = result.outputs
  .reduce((sum, output) => sum + output.size, 0);
const sizeMB = (totalSize / 1024 / 1024).toFixed(2);

console.log(`✓ Built ${result.outputs.length} files to ${outdir} (${sizeMB} MB total)`);

// Watch mode
if (isDev && args.includes("--watch")) {
  console.log("Watching for changes...");
  // Note: Bun doesn't have built-in watch mode for build API yet
  // This would need a file watcher implementation
  console.warn("Watch mode not yet implemented in Bun build script");
}
