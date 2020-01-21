package main

import (
	apiserver "github.com/Andronovdima/tpark-db-forum/internal/app"
	"log"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
