/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.templ",
    "./main.go",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      // Using Catppuccin color palette via @catppuccin/tailwindcss
      // Colors are available as ctp-* classes (e.g., ctp-blue, ctp-surface0, etc.)
    },
  },
  plugins: [require("@catppuccin/tailwindcss")({
      prefix: "ctp",
      // which flavour of colours to use by default, in the `:root`
      defaultFlavour: "mocha",
  }),],
}
