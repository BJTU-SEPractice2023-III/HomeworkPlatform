package main

import (
	"homework_platform/server"
	"log"
)

func main() {
	api := server.InitRouter()

	err := api.Run(":8888")
	if err != nil {
		log.Panicln(err)
	}
}
