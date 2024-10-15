# easel

## Quick Start
### 步骤 1：安装所需工具
#### Protobuf 编译器（protoc）
下载并安装 Protobuf 编译器。安装完成后，确保在终端中运行以下命令可以看到版本号：
```bash
protoc --version
```
#### 安装 Go 插件（protoc-gen-go 和 protoc-gen-go-grpc）：
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
#### 确保 Go 工具链路径已添加到环境变量

### 步骤 2：定义 Protobuf 文件