package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(".env.local", ".env"); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}

	handleConfig()
}

func main() {
	fmt.Println("Microservice v1.0 - Copyright (c) 2022 Anthony Collier")

	handleServer()
}
