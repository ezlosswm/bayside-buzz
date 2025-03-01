/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./cmd/web/**/*.html", "./cmd/web/**/*.templ"],
  theme: {
    extend: {
      backgroundImage: {
        hero: "url('/assets/images/corozal-sign.jpg')",
      },
      boxShadow: {
        "md-x":
          "4px 0 6px -1px rgba(0, 0, 0, 0.05), -4px 0 6px -1px rgba(0, 0, 0, 0.05)",
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    function ({ addUtilities }) {
      addUtilities({
        ".no-scrollbar": {
          "-ms-overflow-style": "none", // IE and Edge
          "scrollbar-width": "none", // Firefox
          "&::-webkit-scrollbar": {
            display: "none", // Chrome, Safari
          },
        },
      });
    },
  ],
};
