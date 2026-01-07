.PHONY: help build test lint clean

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# 后端命令
backend-build: ## 构建后端
	cd backend && go build -o bin/server ./cmd/server

backend-test: ## 运行后端测试
	cd backend && go test ./...

backend-test-coverage: ## 查看后端测试覆盖率（PostID）
	cd backend && go test -coverprofile=coverage.out \
		-coverpkg=./internal/domain/content \
		./test/unit/domain/content/... \
		&& go tool cover -func=coverage.out

backend-test-coverage-html: ## 生成后端测试覆盖率 HTML 报告（PostID）
	cd backend && go test -coverprofile=coverage.out \
		-coverpkg=./internal/domain/content \
		./test/unit/domain/content/... \
		&& go tool cover -html=coverage.out -o coverage.html \
		&& echo "覆盖率报告已生成: backend/coverage.html"

backend-test-coverage-usecase: ## 查看 UseCase 测试覆盖率
	cd backend && go test -coverprofile=coverage.out \
		-coverpkg=./internal/application/content,./internal/application/dto \
		./test/unit/application/content/... \
		&& go tool cover -func=coverage.out

backend-test-coverage-usecase-html: ## 生成 UseCase 测试覆盖率 HTML 报告
	cd backend && go test -coverprofile=coverage.out \
		-coverpkg=./internal/application/content,./internal/application/dto \
		./test/unit/application/content/... \
		&& go tool cover -html=coverage.out -o coverage-usecase.html \
		&& echo "覆盖率报告已生成: backend/coverage-usecase.html"

backend-lint: ## 运行后端代码检查
	cd backend && golangci-lint run

backend-tidy: ## 整理后端依赖
	cd backend && go mod tidy

# 前端命令
frontend-install: ## 安装前端依赖
	cd frontend && npm install

frontend-dev: ## 启动前端开发服务器
	cd frontend && npm run dev

frontend-build: ## 构建前端
	cd frontend && npm run build

frontend-test: ## 运行前端测试
	cd frontend && npm test

frontend-preview: ## 预览前端构建产物
	cd frontend && npm run preview

frontend-verify: ## 验证前端构建（TypeScript 编译 + Vite 构建）
	@echo "🔍 验证前端项目..."
	cd frontend && npm run build && echo "✅ 前端构建验证通过"

# Database seed commands
seed-data: ## 清理并插入12条模拟数据到数据库
	cd backend && go run scripts/seed_data.go

frontend-generate-grpc: ## 生成前端 gRPC Web 代码
	cd frontend && ./scripts/generate-grpc.sh

# 数据库命令
db-migrate-up: ## 运行数据库迁移（向上）
	cd backend && migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up

db-migrate-down: ## 回滚数据库迁移
	cd backend && migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down

# Docker 命令
docker-up: ## 启动 Docker 服务（PostgreSQL + Redis + Backend）
	docker-compose up -d

docker-down: ## 停止 Docker 服务
	docker-compose down

docker-build: ## 构建 Docker 镜像
	docker-compose build

docker-logs: ## 查看 Docker 服务日志
	docker-compose logs -f

docker-ps: ## 查看 Docker 服务状态
	docker-compose ps

docker-restart: ## 重启 Docker 服务
	docker-compose restart

docker-clean: ## 停止并删除所有容器、卷和网络
	docker-compose down -v --remove-orphans

# 测试环境
test-up: ## 启动测试环境
	docker-compose -f docker-compose.test.yml up -d

test-down: ## 停止测试环境
	docker-compose -f docker-compose.test.yml down

test-integration: ## 运行集成测试
	cd backend && go test -v ./test/integration/...

test-integration-repository: ## 运行 Repository 集成测试
	cd backend && go test -v ./test/integration/repository/...

test-integration-cache: ## 运行 Cache 集成测试
	cd backend && go test -v ./test/integration/cache/...

test-integration-usecase: ## 运行 UseCase 集成测试
	cd backend && go test -v ./test/integration/usecase/...

test-unit-usecase: ## 运行 UseCase 单元测试
	cd backend && go test -v ./test/unit/application/content/...

test-unit-usecase-coverage: ## 查看 UseCase 单元测试覆盖率
	cd backend && go test -coverprofile=coverage.out \
		-coverpkg=./internal/application/content,./internal/application/dto \
		./test/unit/application/content/... \
		&& go tool cover -func=coverage.out

test-unit-usecase-coverage-html: ## 生成 UseCase 单元测试覆盖率 HTML 报告
	cd backend && go test -coverprofile=coverage.out \
		-coverpkg=./internal/application/content,./internal/application/dto \
		./test/unit/application/content/... \
		&& go tool cover -html=coverage.out -o coverage-usecase-unit.html \
		&& echo "覆盖率报告已生成: backend/coverage-usecase-unit.html"

test-unit-grpc: ## 运行 gRPC Handler 单元测试
	cd backend && go test -v ./test/unit/presentation/grpc/...

test-unit-grpc-coverage: ## 查看 gRPC Handler 单元测试覆盖率
	cd backend && go test -coverprofile=coverage-grpc.out \
		-coverpkg=./internal/presentation/grpc \
		./test/unit/presentation/grpc/... \
		&& go tool cover -func=coverage-grpc.out

test-unit-grpc-coverage-html: ## 生成 gRPC Handler 单元测试覆盖率 HTML 报告
	cd backend && go test -coverprofile=coverage-grpc.out \
		-coverpkg=./internal/presentation/grpc \
		./test/unit/presentation/grpc/... \
		&& go tool cover -html=coverage-grpc.out -o coverage-grpc.html \
		&& echo "覆盖率报告已生成: backend/coverage-grpc.html"

test-e2e-grpc: ## 运行 gRPC E2E 测试（需要先启动测试环境: make test-up）
	cd backend && go test -v ./test/e2e/scenarios/...

test-integration-debug: ## 运行调试集成测试（不清理数据，可查看数据库）
	cd backend && go test -tags=debug -v ./test/integration/repository/... -run TestPostRepository_Save_Debug

# gRPC 代码生成
generate-proto: ## 生成 gRPC Go 代码
	cd backend && ./scripts/generate.sh

# 清理
clean: ## 清理构建文件
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules

