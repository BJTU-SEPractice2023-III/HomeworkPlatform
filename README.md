# HomeworkPlatform

## 构建

### 1. 安装 statik

本项目使用 [rakyll/statik: Embed files into a Go executable](https://github.com/rakyll/statik) 将静态文件打包入可执行文件中

```shell
go get github.com/rakyll/statik
go install github.com/rakyll/statik
```

### 2. 构建前端

```shell
git submodule update --init
cd assets
pnpm i
pnpm build
```

### 3. 打包前端到 go 模块

```shell
statik -src=assets/build -f
```

### 3. 构建可执行文件

```shell
set GOARCH=amd64
set GOOS=linux
go build -o Builds/v1.0/HomeworkPlatform-1.0.0 homework_platform
set GOOS=windows
go build -o Builds/v1.0/ACH-1.0.0.exe homework_platform
```
