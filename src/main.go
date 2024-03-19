package main

import "github.com/raulito1500/merkadapp/server"

func main() {
	api := server.NewApi()
	api.Run()
}
