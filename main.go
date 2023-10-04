package main

import (
	"embed"
	"homework_platform/internal/bootstrap"
	"homework_platform/server"
	"log"
)


// go1.16后新加入的功能，将文件或目录作为一个文件系统嵌入到二进制文件中
// 这里就是将前端的构建产物目录嵌入到了二进制文件中
// 所以如果这里报错，大概是没有构建前端，目录不存在
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
