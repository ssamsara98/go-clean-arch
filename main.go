package main

import (
	"go-clean-arch/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = bootstrap.RootApp.Execute()
}
