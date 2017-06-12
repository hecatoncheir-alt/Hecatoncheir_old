package main

import (
	"fmt"

	http "github.com/hecatoncheir/Hecatoncheir/http"
)

func main() {
	httpServer := http.NewEngine("v1.0")
	fmt.Println(httpServer.APIVersion)
}
