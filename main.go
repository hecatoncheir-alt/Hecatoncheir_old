package main

import (
	http "github.com/hecatoncheir/Hecatoncheir/http"
	socket "github.com/hecatoncheir/Hecatoncheir/socket"
)

func main() {
	httpServer := http.NewEngine("v1.0")
	socketServer := socket.NewEngine("v1.0")

	httpServer.Router.HandlerFunc("GET", "/", socketServer.ClientConnectedHandler)

	httpServer.PowerUp("0.0.0.0", 8080)
}
