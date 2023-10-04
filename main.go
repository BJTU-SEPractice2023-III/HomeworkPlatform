package main

import (
	"embed"
	"homework_platform/internal/bootstrap"
	"homework_platform/server"
	"log"
    "io/fs"
    "fmt"
)

//go:embed all:assets/dist/public/*
var f embed.FS

func init() {
	bootstrap.InitStatic(f)
    _ = fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
}

func main() {
	api := server.InitRouter()

	err := api.Run(":8888")
	if err != nil {
		log.Panicln(err)
	}
}
