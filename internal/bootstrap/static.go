package bootstrap

import (
	"embed"
	"log"
	"net/http"

	"io/fs"
)

var StaticFS http.FileSystem

func InitStatic(statics embed.FS) {
	log.Println("[bootStrap/InitStaticFS]: Initializing...")
	var err error
	embedFS, err := fs.Sub(statics, "assets/build")
	if err != nil {
		log.Panicf("Failed to initialize static resources: %s", err)
	}

	StaticFS = http.FS(embedFS)
}
