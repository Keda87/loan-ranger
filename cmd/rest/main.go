package main

import "loan-ranger/internal/app/rest"

func main() {
	server := rest.NewHTTPServer()
	server.Start()
	defer server.Stop()
}
