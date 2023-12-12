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
                "foreground-dark": "#0e0e0e",
                "foreground": "#121212",
                "base": "#000000",

                "primary": "#4b7bec",
                "primary-dark": "#375bad",
                "secondary": "#778ca3",
            },
        },
    },
    plugins: [],
}