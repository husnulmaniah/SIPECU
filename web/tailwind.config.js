/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // ── Palet Instansi Pemerintah ──────────────────────────
        // Navy (Primary Dominant) – formal, profesional, tepercaya
        navy: {
          50:  '#EBF0F8',
          100: '#C7D5EB',
          200: '#A0B9DC',
          300: '#6F92C8',
          400: '#4270B4',
          500: '#1A365D', // Base navy
          600: '#152D4E',
          700: '#10233D',
          800: '#0B192C',
          900: '#060F1A',
        },
        // Gold / Mustard (Accent) – kejayaan, CTA, kekayaan budaya
        gold: {
          50:  '#FEFCE8',
          100: '#FEF3C7',
          200: '#FDE68A',
          300: '#FCD34D',
          400: '#FBBF24',
          500: '#D69E2E', // Base gold
          600: '#B7851F',
          700: '#926B14',
          800: '#6D500D',
          900: '#4A3608',
        },
        // Off-White (Background) – latar ramah mata
        offwhite: '#F7FAFC',
        // Dark Gray (Text) – teks utama
        charcoal: {
          50:  '#F7F8FA',
          100: '#EDF0F4',
          200: '#D4DAE4',
          300: '#B0BAC9',
          400: '#8494A8',
          500: '#2D3748', // Base dark gray / text
          600: '#252F3B',
          700: '#1C252E',
          800: '#141A20',
          900: '#0C0F14',
        },
        // Semantic alias – tetap digunakan di seluruh komponen Vue via class primary-*
        primary: {
          50:  '#EBF0F8',
          100: '#C7D5EB',
          200: '#A0B9DC',
          300: '#6F92C8',
          400: '#4270B4',
          500: '#1A365D',
          600: '#152D4E',
          700: '#10233D',
          800: '#0B192C',
          900: '#060F1A',
        },
        accent: {
          50:  '#FEFCE8',
          100: '#FEF3C7',
          300: '#FCD34D',
          400: '#FBBF24',
          500: '#D69E2E',
          600: '#B7851F',
          700: '#926B14',
        },
        // Background card / panel
        darkbg:   '#F7FAFC',   // Off-white page background
        darkcard: '#EDF2F7',   // Slightly darker card
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif']
      }
    },
  },
  plugins: [],
}
