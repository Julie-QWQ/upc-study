package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/study-upc/backend/internal/pkg/config"
	"github.com/study-upc/backend/internal/pkg/logger"
	"go.uber.org/zap"
)

// initSchemaMigrations 创建迁移版本跟踪表
func initSchemaMigrations() error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	return DB.Exec(sql).Error
}

// isMigrationApplied 检查迁移是否已应用
func isMigrationApplied(version string) (bool, error) {
	var count int64
	err := DB.Table("schema_migrations").Where("version = ?", version).Count(&count).Error
	return count > 0, err
}

// recordMigration 记录已执行的迁移
func recordMigration(version string) error {
	return DB.Table("schema_migrations").Create(map[string]interface{}{
		"version": version,
	}).Error
}

// extractMigrationVersion 从文件名提取迁移版本号
// 例如: "001_init_schema.up.sql" -> "001"
func extractMigrationVersion(filename string) string {
	base := filepath.Base(filename)
	// 去掉 .up.sql 或 .down.sql 后缀
	version := strings.TrimSuffix(base, ".up.sql")
	version = strings.TrimSuffix(version, ".down.sql")
	return version
}

func findMigrationFiles() (string, []string, error) {
	candidates := []string{"migrations", filepath.Join("backend", "migrations")}

	var searchPaths []string
	if cwd, err := os.Getwd(); err == nil {
		for _, dir := range candidates {
			searchPaths = append(searchPaths, filepath.Join(cwd, dir))
		}
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		for _, dir := range candidates {
			searchPaths = append(searchPaths, filepath.Join(exeDir, dir))
		}
	}
	searchPaths = append(searchPaths, candidates...)

	for _, dir := range searchPaths {
		files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
		if err != nil {
			continue
		}
		if len(files) > 0 {
			return dir, files, nil
		}
	}

	return "", nil, fmt.Errorf("no migration files found in %v", searchPaths)
}

// RunMigrations 执行数据库迁移
func RunMigrations(cfg *config.Config) error {
	// 初始化迁移版本跟踪表
	if err := initSchemaMigrations(); err != nil {
		return fmt.Errorf("初始化迁移跟踪表失败: %w", err)
	}

	// 获取所有 .up.sql 迁移文件
	_, files, err := findMigrationFiles()
	if err != nil {
		return fmt.Errorf("查找迁移文件失败: %w", err)
	}

	// 按文件名排序执行
	sort.Strings(files)

	appliedCount := 0
	skippedCount := 0

	for _, file := range files {
		version := extractMigrationVersion(file)

		// 检查迁移是否已应用
		applied, err := isMigrationApplied(version)
		if err != nil {
			return fmt.Errorf("检查迁移版本失败 %s: %w", file, err)
		}

		if applied {
			logger.Debug("跳过已应用的迁移", zap.String("file", file), zap.String("version", version))
			skippedCount++
			continue
		}

		logger.Info("执行迁移", zap.String("file", file), zap.String("version", version))

		// 读取迁移 SQL 文件
		upSQL, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("读取迁移文件失败 %s: %w", file, err)
		}

		// 执行迁移
		if err := DB.Exec(string(upSQL)).Error; err != nil {
			// 记录错误但继续执行(可能字段已存在)
			logger.Warn("迁移执行失败", zap.String("file", file), zap.String("version", version), zap.Error(err))
		} else {
			// 记录已执行的迁移
			if err := recordMigration(version); err != nil {
				logger.Warn("记录迁移版本失败", zap.String("version", version), zap.Error(err))
			}
			logger.Info("迁移执行成功", zap.String("file", file), zap.String("version", version))
			appliedCount++
		}
	}

	logger.Info("数据库迁移完成",
		zap.Int("应用数量", appliedCount),
		zap.Int("跳过数量", skippedCount))
	return nil
}

// RollbackMigrations 回滚最后一个迁移
func RollbackMigrations(cfg *config.Config) error {
	// 获取最新应用的迁移版本
	var migration struct {
		Version string
	}
	err := DB.Table("schema_migrations").Order("applied_at DESC").First(&migration).Error
	if err != nil {
		return fmt.Errorf("获取最新迁移版本失败: %w", err)
	}

	// 查找对应的 down 文件
	migrationsDir, _, err := findMigrationFiles()
	if err != nil {
		return fmt.Errorf("查找迁移文件失败: %w", err)
	}
	downFile := filepath.Join(migrationsDir, migration.Version+".down.sql")

	if _, err := os.Stat(downFile); os.IsNotExist(err) {
		return fmt.Errorf("回滚文件不存在: %s", downFile)
	}

	logger.Info("执行回滚", zap.String("version", migration.Version), zap.String("file", downFile))

	// 读取回滚 SQL 文件
	downSQL, err := os.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("读取回滚文件失败 %s: %w", downFile, err)
	}

	// 执行回滚
	if err := DB.Exec(string(downSQL)).Error; err != nil {
		return fmt.Errorf("回滚执行失败: %w", err)
	}

	// 删除迁移记录
	if err := DB.Table("schema_migrations").Where("version = ?", migration.Version).Delete(nil).Error; err != nil {
		return fmt.Errorf("删除迁移记录失败: %w", err)
	}

	logger.Info("回滚执行成功", zap.String("version", migration.Version))
	return nil
}
