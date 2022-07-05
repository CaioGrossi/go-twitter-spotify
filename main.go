package main

import (
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	client := http.Client{}
	connectTwitterHttpStream(client)
}
