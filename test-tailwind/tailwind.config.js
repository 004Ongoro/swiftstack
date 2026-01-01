import type { Config } from "tailwindcss";

const config: Config = {
  // We add 'class' mode here to fix the dark mode issue properly
  darkMode: 'class',
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // You can add custom colors here later if needed
      },
    },
  },
  plugins: [],
};
export default config;