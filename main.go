package main

import (
	"fmt"
	"gopherreq/gopherreq"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	address, hasEnvAddress := os.LookupEnv("HTTP_DOMAIN")
	if !hasEnvAddress {
		address = "localhost:8811"
	}

	config := gopherreq.Config{
		Domain:  address,
		Timeout: 4000,
	}

	server, err := gopherreq.NewServer(config)

	fmt.Printf("Server started listening and is accepting connections on the fly.")

	if err != nil {
		fmt.Printf("Error ocurred while listening on server - %v", err)
		os.Exit(1)
	}

	server.Listen()

	defer server.ShutDown()
}
