import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        forest: {
          900: '#0F332E',
          850: '#143F3A',
          800: '#184F48',
          700: '#20635B',
          600: '#2B7A70',
        },
        mint: {
          50: '#F4FAF7',
          100: '#E7F5EF',
          200: '#CBEBDD',
          300: '#A4DCC5',
          400: '#64BE9F',
          500: '#38A783',
          600: '#268767',
        },
        surface: {
          light: '#EAF5F0',
          card: 'rgba(255, 255, 255, 0.72)',
          border: 'rgba(255, 255, 255, 0.65)',
        },
      },
      fontFamily: {
        sans: ['"Plus Jakarta Sans"', 'Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'glass': '0 8px 32px 0 rgba(24, 79, 72, 0.05)',
        'card-subtle': '0 4px 20px -2px rgba(24, 79, 72, 0.06)',
      },
    },
  },
  plugins: [],
};
export default config;
