package main

import (
	"fmt"
	httpproto "http-v1_1/http-proto"
	"os"
)

func main() {

	config := httpproto.Config{
		Domain:  "localhost:8811",
		Timeout: 4000,
	}

	server, err := httpproto.NewServer(config)

	fmt.Printf("Server started listening and is accepting connections on the fly.")

	if err != nil {
		fmt.Printf("Error ocurred while listening on server - %v", err)
		os.Exit(1)
	}

	server.Listen()

	defer server.ShutDown()
}
