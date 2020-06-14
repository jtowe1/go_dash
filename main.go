package main

import (
    "github.com/joho/godotenv"
    "jeremiahtowe.com/go_dash/goDash"
    "log"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    goDash.Run()
}
