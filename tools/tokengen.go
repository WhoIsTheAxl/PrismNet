package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sandertv/gophertunnel/minecraft/auth"
)

func main() {
	fmt.Println("Requesting Microsoft login...")

	token, err := auth.RequestLiveToken()
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("token.json", data, 0600); err != nil {
		panic(err)
	}

	fmt.Println("token.json created successfully.")
}