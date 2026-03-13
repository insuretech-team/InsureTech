import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  output: "standalone",

  // Transpile the local-tgz SDK so Next.js processes its ESM output correctly.
  transpilePackages: ["@lifeplus/insuretech-sdk"],

  // Webpack alias (used by `next build` and `next start`)
  webpack(config) {
    config.resolve.alias = {
      ...config.resolve.alias,
      "@lifeplus/insuretech-sdk": path.resolve(
        __dirname,
        "node_modules/@lifeplus/insuretech-sdk"
      ),
    };
    return config;
  },

  // Turbopack alias must use forward-slash relative path (Windows absolute paths
  // are not supported by Turbopack — "windows imports are not implemented yet")
  turbopack: {
    resolveAlias: {
      "@lifeplus/insuretech-sdk": "./node_modules/@lifeplus/insuretech-sdk/dist/index.mjs",
    },
  },
};

export default nextConfig;
