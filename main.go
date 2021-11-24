package main

import "github.com/mariajdab/post-api-rest/api"

func main() {
	server := api.NewServer()
	server.Run()
}
