package main

import (
	"github.com/joho/godotenv"
	"github.com/ssamsara98/go-clean-arch/src/bootstrap"
)

func main() {
	_ = godotenv.Load()
	_ = bootstrap.RootApp.Execute()
}
