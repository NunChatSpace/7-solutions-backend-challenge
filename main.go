package main

import "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http"

func main() {
	server := http.NewServer()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
