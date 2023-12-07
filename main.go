package main

import (
	"fmt"
	"os"
	"spotisong/api"
	"strings"

	"github.com/joho/godotenv"
)

func Watch(args [] string) {
	tailwind := api.TailWind {
		Version: os.Getenv("TAILWIND_VERSION"),
		Binary: "./.tailwind/",
	}

	err := tailwind.Watch(
		"./app/style.css",
		"./app/static/" + os.Getenv("TAILWIND_OUTPUT"),
	)

	if err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load '.env' file, maybe its missing?")
	}


	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(ACTIONS_RESPONSE)
		return
	}

	if action, ok := ACTIONS[strings.ToLower(os.Args[1])]; ok {
		action(args[1:])
	} else {
		fmt.Println(ACTIONS_RESPONSE)
	}
}

var ACTIONS = map [string] func(args [] string) {
	"watch": Watch,
}

// this is the response that gets 
// printed onto the screen if the
// user provided invalid launch args
const ACTIONS_RESPONSE = "<watch>"