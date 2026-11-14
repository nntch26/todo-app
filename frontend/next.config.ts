import type { NextConfig } from "next";

// Where the Go API runs. Local dev: localhost:8080. Docker (Step 9): http://backend:8080
const apiUrl = process.env.API_URL ?? "http://localhost:8080";

// Fail at startup if API_URL is not a valid URL
new URL(apiUrl);

const nextConfig: NextConfig = {
  /* config options here */
  cacheComponents: true,
  partialPrefetching: true,
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },

  // Proxy /api/* to the Go API, so the browser sees one origin (localhost:3000)
  // and the httpOnly login cookie works without CORS
  rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${apiUrl}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
