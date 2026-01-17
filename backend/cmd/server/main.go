package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/pkg/config"
	"github.com/study-upc/backend/internal/pkg/database"
	"github.com/study-upc/backend/internal/pkg/logger"
	"github.com/study-upc/backend/internal/router"
	"go.uber.org/zap"

	_ "github.com/study-upc/backend/docs" // swagger docs
)

// @title           Study-UPC API
// @version         1.0
// @description     学习资料托管平台 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.study-upc.com/support
// @contact.email  support@study-upc.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		appEnv := os.Getenv("APP_ENV")
		ginMode := os.Getenv("GIN_MODE")
		configPath = pickConfigPath(appEnv, ginMode)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("启动 Study-UPC 后端服务")

	// 初始化数据库
	if err := database.InitPostgres(cfg); err != nil {
		logger.Fatal("初始化数据库失败", zap.Error(err))
	}
	defer database.ClosePostgres()

	// 执行数据库迁移
	if err := database.RunMigrations(cfg); err != nil {
		logger.Warn("数据库迁移执行失败", zap.Error(err))
	}

	// 初始化 Redis
	if err := database.InitRedis(cfg); err != nil {
		logger.Fatal("初始化Redis失败", zap.Error(err))
	}
	defer database.CloseRedis()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 设置路由
	r := router.SetupRouter(cfg)

	// HTTP 服务器
	srv := &http.Server{
		Addr:         cfg.Server.GetAddr(),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("服务器启动",
			zap.String("addr", cfg.Server.GetAddr()),
			zap.String("mode", cfg.Server.Mode),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	logger.Info("服务器已退出")
}

func pickConfigPath(appEnv, ginMode string) string {
	var candidates []string
	if appEnv == "production" || ginMode == "release" {
		candidates = []string{"configs/config.prod.local.yaml", "configs/config.prod.yaml"}
	} else {
		candidates = []string{"configs/config.dev.local.yaml", "configs/config.dev.yaml"}
	}

	paths := resolvePaths(candidates)
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return candidates[len(candidates)-1]
}

func resolvePaths(paths []string) []string {
	var resolved []string
	// Check current working directory first.
	if cwd, err := os.Getwd(); err == nil {
		for _, p := range paths {
			resolved = append(resolved, filepath.Join(cwd, p))
		}
	}
	// Then check relative to the executable.
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		for _, p := range paths {
			resolved = append(resolved, filepath.Join(exeDir, p))
		}
	}
	// Finally, fall back to raw relative paths.
	resolved = append(resolved, paths...)
	return resolved
}
