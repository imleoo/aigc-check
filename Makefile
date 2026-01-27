.PHONY: build test coverage clean install lint help run

# 变量定义
BINARY_NAME=aigc-check
VERSION=0.1.0
BUILD_DIR=bin
MAIN_PATH=cmd/aigc-check/main.go

# 默认目标
.DEFAULT_GOAL := help

# 构建项目
build: ## 编译项目
	@echo "正在编译 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行测试
test: ## 运行测试
	@echo "正在运行测试..."
	@go test -v ./...

# 测试覆盖率
coverage: ## 生成测试覆盖率报告
	@echo "正在生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 清理构建产物
clean: ## 清理构建产物
	@echo "正在清理..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "清理完成"

# 安装到系统
install: build ## 安装到系统
	@echo "正在安装 $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "安装完成: /usr/local/bin/$(BINARY_NAME)"

# 代码检查
lint: ## 运行代码检查
	@echo "正在运行代码检查..."
	@go vet ./...
	@go fmt ./...
	@echo "代码检查完成"

# 运行示例
run: build ## 运行示例
	@echo "运行示例检测..."
	@$(BUILD_DIR)/$(BINARY_NAME) -f test/testdata/sample.txt

# 依赖管理
deps: ## 下载依赖
	@echo "正在下载依赖..."
	@go mod download
	@go mod tidy
	@echo "依赖下载完成"

# 跨平台编译
build-all: ## 跨平台编译
	@echo "正在进行跨平台编译..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "跨平台编译完成"

# 帮助信息
help: ## 显示帮助信息
	@echo "AIGC-Check Makefile"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
