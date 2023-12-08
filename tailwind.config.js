/** @type {import('tailwindcss').Config} */

module.exports = {
    darkMode: 'class',
    mode: "jit",
    content: [
        "./app/**/*.{html,js}",
    ],
    theme: {
        extend: {
            colors: {
                "primary": "#121212",
                "secondary": "#000000",
            },
        },
    },
    plugins: [],
}