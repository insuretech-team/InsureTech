import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  output: "standalone",
  transpilePackages: ["@lifeplus/insuretech-sdk"],
  webpack(config) {
    config.resolve.alias = {
      ...config.resolve.alias,
      "@lifeplus/insuretech-sdk": path.resolve(
        __dirname,
        "node_modules/@lifeplus/insuretech-sdk",
      ),
    };

    return config;
  },
  turbopack: {
    resolveAlias: {
      "@lifeplus/insuretech-sdk": "./node_modules/@lifeplus/insuretech-sdk/dist/index.mjs",
    },
  },
};

export default nextConfig;
