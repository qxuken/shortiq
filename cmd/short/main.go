package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/qxuken/short/internal/app"
	"github.com/qxuken/short/internal/auth"
)

func hashCmd() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments: no token provided")
		os.Exit(1)
	}
	token := os.Args[2]
	phc, err := auth.GeneratePHCHash([]byte(token))
	if err != nil {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	fmt.Println(phc)
}

func main() {
	godotenv.Load()

	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "hash":
		hashCmd()
	default:
		app.RunApp()
	}

}
