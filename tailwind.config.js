/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/template/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
};
