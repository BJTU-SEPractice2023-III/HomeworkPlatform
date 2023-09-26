package main

import (
	"embed"
	"homework_platform/internal/bootstrap"
	"homework_platform/server"
	"log"
)

//go:embed assets/build/*
var f embed.FS

func init() {
	bootstrap.InitStatic(f)
}

func main() {
	api := server.InitRouter()

	err := api.Run(":8888")
	if err != nil {
		log.Panicln(err)
	}
}
