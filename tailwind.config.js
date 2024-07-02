/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'selector',
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", "./**/*.svg", ],
  theme: {
    extend: {},
  },
  plugins: [
    require("@catppuccin/tailwindcss")({
      defaultFlavour: "latte",
  }),
  ],
}


