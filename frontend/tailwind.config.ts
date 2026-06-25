import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        ink: "#17211d",
        road: "#2f6153",
        signal: "#f4b63d",
        clay: "#d65f45",
        mist: "#eef5f1",
      },
    },
  },
  plugins: [],
};

export default config;
