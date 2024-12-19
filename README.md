# lmkproxy

lmkproxy 是一个基于 Go 语言的网络代理项目，包含客户端和服务器端。

## 目录结构

```txt
go/ 
├── src/ 
│ ├── client/ 
│ │ └── client.go 
│ ├── server/ 
│ │ └── server.go 
│ ├── internal/ 
│ │ └── core/ 
│ │ ├── cypher.go 
│ │ └── trans.go 
├── configs/ 
│ └── config.txt 
├── go.mod 
├── go.sum 
└── README.md
```

## 安装

1. 克隆仓库：

    ```sh
    git clone https://github.com/LLLMMKK/lmkproxy.git
    cd lmkproxy
    ```

2. 安装依赖项：

    ```sh
    go mod download
    ```

## 配置

在 `configs/` 目录下创建 `config.txt` 文件，并添加以下内容：

```txt
vps_address=YOUR_VPS_ADDRESS
vps_port=YOUR_VPS_PORT
```

## 使用

### 启动服务器

1. 进入 `src/server` 目录：

    ```sh
    cd src/server
    ```

2. 构建并运行服务器：

    ```sh
    go build -o server
    ./server
    ```

### 启动客户端

1. 进入 `src/client` 目录：

    ```sh
    cd src/client
    ```

2. 构建并运行客户端：

    ```sh
    go build -o client
    ./client
    ```

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 仓库
2. 创建新分支 (`git checkout -b feature-branch`)
3. 提交更改 (`git commit -am 'Add new feature'`)
4. 推送到分支 (`git push origin feature-branch`)
5. 创建 Pull Request

## 许可证

本项目使用 MIT 许可证。