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

                "primary": "#ff6b81",
                "primary-dark": "#ff4757",
                "secondary": "#3742fa", 
                "secondary-dark": "#282fad", 
            },
        },
    },
    plugins: [],
}