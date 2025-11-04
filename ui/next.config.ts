import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: 'standalone',

  // Allow external API calls
  async rewrites() {
    return [
      {
        source: '/api/server/:path*',
        destination: 'http://go-app:8080/api/:path*',
      },
    ];
  },
};

export default nextConfig;
