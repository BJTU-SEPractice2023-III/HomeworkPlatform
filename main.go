package main

import (
	"embed"
	// "fmt"
	// "io/fs"
	"homework_platform/internal/bootstrap"
	"homework_platform/server"
	"log"
	"os"
	"os/exec"
)

// go1.16后新加入的功能，将文件或目录作为一个文件系统嵌入到二进制文件中
// 这里就是将前端的构建产物目录嵌入到了二进制文件中
// 所以如果这里报错，大概是没有构建前端，目录不存在
//
//go:embed all:assets/dist/*
var f embed.FS

func init() {
	bootstrap.InitStatic(f)

	// _ = fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println(path)
	// 	return nil
	// })
}

func runPnpmDev() {
	projectPath := "assets"
	err := os.Chdir(projectPath)
	if err != nil {
		log.Fatal(err)
	}

	// 执行 PNPM dev 命令
	cmd := exec.Command("pnpm", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// 在这里执行其他 Golang 操作，例如启动 Golang 服务器

	// 等待 PNPM dev 进程结束
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if bootstrap.Dev {
		go runPnpmDev()
	}
	api := server.InitRouter()

	err := api.Run(":8888")
	if err != nil {
		log.Panicln(err)
	}
}
