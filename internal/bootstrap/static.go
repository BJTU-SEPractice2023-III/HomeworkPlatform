package bootstrap

import (
	"fmt"
	"log"
	"net/http"

	"io/fs"
)

var StaticFS http.FileSystem

func InitStatic(staticZip string) {
	log.Println("[bootStrap/InitStaticFS]: Initializing...")

	statics := NewFS(staticZip)

	var err error
	embedFS, err := fs.Sub(statics, "assets/dist/public")
	if err != nil {
		log.Panicf("Failed to initialize static resources: %s", err)
	}

	StaticFS = http.FS(embedFS)

    _ = fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
}
