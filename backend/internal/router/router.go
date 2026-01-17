package router

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/study-upc/backend/internal/handler"
	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/pkg/config"
	"github.com/study-upc/backend/internal/pkg/database"
	"github.com/study-upc/backend/internal/pkg/logger"
	"github.com/study-upc/backend/internal/pkg/oss"
	"github.com/study-upc/backend/internal/pkg/utils"
	"github.com/study-upc/backend/internal/repository"
	"github.com/study-upc/backend/internal/service"
	"go.uber.org/zap"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// 健康检查
	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)
	r.GET("/liveness", healthHandler.Liveness)

	// Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化依赖
	db := database.DB
	redisClient := database.RDB

	// 初始化 JWT 管理器
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.ExpireTime)*time.Hour,
		time.Duration(cfg.JWT.ExpireTime)*24*time.Hour, // refresh token 有效期为 access token 的 24 倍
		"study-upc",
	)

	// 初始化 Repository 层
	userRepo := repository.NewUserRepository(db)
	materialRepo := repository.NewMaterialRepository(db)
	materialCategoryRepo := repository.NewMaterialCategoryRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	downloadRepo := repository.NewDownloadRecordRepository(db)
	reportRepo := repository.NewReportRepository(db)
	committeeRepo := repository.NewCommitteeRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	searchHistoryRepo := repository.NewSearchHistoryRepository(db)
	hotKeywordRepo := repository.NewHotKeywordRepository(db)
	statisticsRepo := repository.NewStatisticsRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)

	// 初始化 OSS 服务
	ossClient, err := oss.NewMinIOClient(&oss.MinIOConfig{
		Endpoint:  cfg.OSS.Endpoint,
		AccessKey: cfg.OSS.AccessKey,
		SecretKey: cfg.OSS.SecretKey,
		Bucket:    cfg.OSS.BucketName,
		Region:    cfg.OSS.Region,
		UseSSL:    cfg.OSS.UseSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("初始化 OSS 客户端失败: %v", err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ossClient.TestConnection(ctx); err != nil {
		logger.Error("OSS 连接测试失败", zap.Error(err), zap.String("endpoint", cfg.OSS.Endpoint), zap.String("bucket", cfg.OSS.BucketName))
		panic(fmt.Sprintf("OSS 连接测试失败: %v", err))
	}
	logger.Info("OSS 连接成功", zap.String("endpoint", cfg.OSS.Endpoint), zap.String("bucket", cfg.OSS.BucketName), zap.Bool("ssl", cfg.OSS.UseSSL))
	// 最大文件大小 512MB，上传签名 1 小时有效期，下载签名 24 小时有效期
	ossService := oss.NewOSSService(ossClient, 536870912, 1*time.Hour, 24*time.Hour)

	// 初始化 Service 层
	authService := service.NewAuthService(userRepo, jwtManager, redisClient)
	materialService := service.NewMaterialService(materialRepo, favoriteRepo, downloadRepo, materialCategoryRepo, adminRepo, ossService, redisClient)
	materialCategoryService := service.NewMaterialCategoryService(materialCategoryRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, materialRepo)
	reportService := service.NewReportService(reportRepo, materialRepo)
	committeeService := service.NewCommitteeService(committeeRepo, userRepo, reviewRepo)
	reviewService := service.NewReviewService(materialRepo, committeeRepo, reportRepo, reviewRepo, userRepo)
	notificationService := service.NewNotificationService(notificationRepo, userRepo)
	searchService := service.NewSearchService(db, materialRepo, searchHistoryRepo, hotKeywordRepo, downloadRepo)
	recommendationService := service.NewRecommendationService(db, materialRepo, downloadRepo, favoriteRepo)
	statisticsService := service.NewStatisticsService(statisticsRepo)
	adminService := service.NewAdminService(adminRepo, userRepo, materialRepo)
	announcementService := service.NewAnnouncementService(announcementRepo, userRepo)

	// 访问日志中间件(需要在 statisticsService 初始化后注册)
	r.Use(middleware.AccessLog(statisticsService))

	// 设置通知服务以解决循环依赖
	committeeService.SetNotificationService(notificationService)
	reviewService.SetNotificationService(notificationService)

	// 初始化 Handler 层
	authHandler := handler.NewAuthHandler(authService, statisticsService)
	materialHandler := handler.NewMaterialHandler(materialService, favoriteService, reportService, downloadRepo)
	materialCategoryHandler := handler.NewMaterialCategoryHandler(materialCategoryService)
	committeeHandler := handler.NewCommitteeHandler(committeeService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	searchHandler := handler.NewSearchHandler(searchService, recommendationService)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)
	adminHandler := handler.NewAdminHandler(adminService)
	announcementHandler := handler.NewAnnouncementHandler(announcementService)
	systemHandler := handler.NewSystemHandler(adminService)

	// 从数据库加载系统配置并应用到 OSS 服务
	if uploadConfig, err := adminService.GetSystemConfig("allowed_file_types"); err == nil && uploadConfig.ConfigValue != "" {
		// 解析允许的文件类型
		typesStr := uploadConfig.ConfigValue
		if typesStr != "" {
			// 手动分割字符串
			allowedTypes := make([]string, 0)
			for _, t := range []string{"pdf", "docx", "doc", "pptx", "ppt", "txt", "md", "zip", "rar"} {
				// 检查是否在配置中
				for _, cfg := range strings.Split(typesStr, ",") {
					if strings.TrimSpace(cfg) == t {
						allowedTypes = append(allowedTypes, t)
						break
					}
				}
			}
			// 如果有配置,更新 OSS 服务的验证器
			if len(allowedTypes) > 0 {
				ossService.UpdateConfig(536870912, allowedTypes)
			}
		}
	}

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		public := v1.Group("")
		{
			// 认证相关
			auth := public.Group("/auth")
			{
				// 登录限流中间件：每个 IP 每小时最多 20 次，每个用户名每 15 分钟最多 5 次
				auth.POST("/login", middleware.LoginRateLimit(redisClient, 20, 5), authHandler.Login)
				auth.POST("/register", authHandler.Register)
				auth.POST("/refresh", authHandler.RefreshToken)
			}

			// 系统配置（公开）
			system := public.Group("/system")
			{
				system.GET("/configs/:key", systemHandler.GetPublicSystemConfig)
				system.GET("/configs", systemHandler.GetPublicSystemConfigs)
				system.GET("/upload-config", systemHandler.GetUploadConfig)
			}
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(jwtManager, redisClient))
		{
			// 资料类型相关（所有认证用户可访问）
			materialCategories := protected.Group("/material-categories")
			{
				materialCategories.GET("", materialCategoryHandler.List)               // 获取资料类型列表
				materialCategories.GET("/:id", materialCategoryHandler.GetByID)        // 获取资料类型详情
			}

			// 资料类型管理（管理员权限）
			adminMaterialCategories := protected.Group("/admin/material-categories")
			adminMaterialCategories.Use(middleware.RequireAdmin())
			{
				adminMaterialCategories.POST("", materialCategoryHandler.Create)               // 创建资料类型
				adminMaterialCategories.PUT("/:id", materialCategoryHandler.Update)            // 更新资料类型
				adminMaterialCategories.DELETE("/:id", materialCategoryHandler.Delete)         // 删除资料类型
				adminMaterialCategories.POST("/:id/toggle", materialCategoryHandler.ToggleStatus) // 切换启用状态
			}

			// 公告相关（所有认证用户可访问）
			announcements := protected.Group("/announcements")
			{
				announcements.GET("/active", announcementHandler.GetActiveAnnouncements) // 获取启用的公告
				announcements.GET("/:id", announcementHandler.GetAnnouncement)         // 公告详情
			}

			// 公告管理（管理员权限）
			adminAnnouncements := protected.Group("/announcements")
			adminAnnouncements.Use(middleware.RequireAdmin())
			{
				adminAnnouncements.GET("", announcementHandler.ListAnnouncements)           // 公告列表
				adminAnnouncements.POST("", announcementHandler.CreateAnnouncement)         // 创建公告
				adminAnnouncements.PUT("/:id", announcementHandler.UpdateAnnouncement)      // 更新公告
				adminAnnouncements.DELETE("/:id", announcementHandler.DeleteAnnouncement)   // 删除公告
			}
			// 认证相关
			auth := protected.Group("/auth")
			{
				auth.POST("/logout", authHandler.Logout)
				auth.POST("/change-password", authHandler.ChangePassword)
				auth.GET("/me", authHandler.GetUserInfo)
			}

			// 资料管理路由
			materials := protected.Group("/materials")
			{
				// 公开接口（所有认证用户可访问）
				materials.GET("", materialHandler.ListMaterials)           // 资料列表
				materials.GET("/search", materialHandler.SearchMaterials) // 搜索资料
				materials.GET("/:id", materialHandler.GetMaterial)        // 资料详情

				// 下载接口（所有认证用户可访问）
				materials.GET("/:id/download", materialHandler.GetDownloadURL) // 获取下载链接

				// 收藏相关（所有认证用户）
				materials.POST("/:id/favorite", materialHandler.AddFavorite)      // 添加收藏
				materials.DELETE("/:id/favorite", materialHandler.RemoveFavorite) // 取消收藏

				// 举报相关（所有认证用户）
				materials.POST("/:id/report", materialHandler.CreateReport) // 创建举报

				// 学委及以上权限
				committee := materials.Use(middleware.RequireCommittee())
				{
					committee.POST("", materialHandler.CreateMaterial)               // 创建资料
					committee.PUT("/:id", materialHandler.UpdateMaterial)            // 更新资料
					committee.POST("/upload-signature", materialHandler.GetUploadSignature) // 获取上传签名
					committee.POST("/delete-uploaded-file", materialHandler.DeleteUploadedFile) // 删除已上传文件
				}

				// 管理员权限
				admin := materials.Use(middleware.RequireAdmin())
				{
					admin.DELETE("/:id", materialHandler.DeleteMaterial)    // 删除资料
					admin.POST("/:id/review", materialHandler.ReviewMaterial) // 审核资料
				}
			}

			// 收藏列表
			protected.GET("/favorites", materialHandler.ListFavorites)

			// 下载记录列表
			protected.GET("/downloads", materialHandler.ListDownloadRecords)

			// 举报管理（管理员）
			adminReports := protected.Group("/admin/reports")
			adminReports.Use(middleware.RequireAdmin())
			{
				adminReports.GET("", materialHandler.ListReports)             // 举报列表
				adminReports.GET("/:id", materialHandler.GetReport)           // 举报详情
				adminReports.POST("/:id/handle", reviewHandler.HandleReport) // 处理举报
			}

			// 学委申请相关
			user := protected.Group("/user")
			{
				// 申请学委
				user.POST("/apply-committee", committeeHandler.ApplyForCommittee)
				// 我的申请列表
				user.GET("/applications", committeeHandler.ListMyApplications)
				// 申请详情
				user.GET("/applications/:id", committeeHandler.GetApplication)
				// 取消申请
				user.POST("/applications/:id/cancel", committeeHandler.CancelApplication)
			}

			// 管理员-学委申请管理
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				// 统计管理
				statistics := admin.Group("/statistics")
				{
					statistics.GET("/overview", statisticsHandler.GetOverviewStatistics)      // 概览统计
					statistics.GET("/users", statisticsHandler.GetUserStatistics)            // 用户统计
					statistics.GET("/users/trend", statisticsHandler.GetUserTrend)            // 用户趋势
					statistics.GET("/materials", statisticsHandler.GetMaterialStatistics)     // 资料统计
					statistics.GET("/materials/trend", statisticsHandler.GetMaterialTrend)    // 资料趋势
					statistics.GET("/downloads", statisticsHandler.GetDownloadStatistics)      // 下载统计
					statistics.GET("/downloads/trend", statisticsHandler.GetDownloadTrend)     // 下载趋势
					statistics.GET("/applications", statisticsHandler.GetApplicationStatistics) // 申请统计
					statistics.GET("/visits", statisticsHandler.GetVisitStatistics)            // 访问统计
					statistics.GET("/visits/trend", statisticsHandler.GetVisitTrend)           // 访问趋势
				}

				// 用户管理
				users := admin.Group("/users")
				{
					users.GET("", adminHandler.ListUsers)                              // 用户列表
					users.GET("/:id", adminHandler.GetUserDetail)                       // 用户详情
					users.PUT("/:id", adminHandler.UpdateUserInfo)                      // 更新用户信息
					users.PUT("/:id/status", adminHandler.UpdateUserStatus)             // 更新用户状态
					users.DELETE("/:id", adminHandler.DeleteUser)                       // 删除用户
				}

				// 系统配置管理
				configs := admin.Group("/configs")
				{
					configs.GET("", adminHandler.ListSystemConfigs)                    // 配置列表
					configs.POST("", adminHandler.CreateSystemConfig)                  // 创建配置
					configs.GET("/:key", adminHandler.GetSystemConfig)                 // 获取配置
					configs.PUT("", adminHandler.UpdateSystemConfig)                   // 更新配置
					configs.DELETE("/:key", adminHandler.DeleteSystemConfig)           // 删除配置
				}

				// 学委申请列表
				admin.GET("/applications", committeeHandler.ListApplications)
				// 审核学委申请
				admin.POST("/applications/:id/review", committeeHandler.ReviewApplication)
				// 待审核申请数量
				admin.GET("/applications/pending/count", committeeHandler.GetPendingCount)

				// 待审核资料列表
				admin.GET("/materials/pending", materialHandler.ListPendingMaterials)
				// 已审核资料列表
				admin.GET("/materials/reviewed", materialHandler.ListReviewedMaterials)
				// 审核资料
				admin.POST("/materials/:id/review", reviewHandler.ReviewMaterial)

				// 审核历史
				admin.GET("/review/history", reviewHandler.GetReviewHistory)
				// 审核人统计
				admin.GET("/reviewers/:id/statistics", reviewHandler.GetReviewerStatistics)
			}

			// 通知相关
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.ListNotifications)            // 通知列表
				notifications.GET("/unread", notificationHandler.GetUnreadNotifications) // 未读通知
				notifications.GET("/unread/count", notificationHandler.GetUnreadCount)   // 未读数量
				notifications.POST("/:id/read", notificationHandler.MarkAsRead)          // 标记已读
				notifications.POST("/read-all", notificationHandler.MarkAllAsRead)       // 全部标记已读
				notifications.DELETE("/:id", notificationHandler.DeleteNotification)     // 删除通知
			}

			// 搜索和推荐相关
			search := protected.Group("/search")
			{
				search.GET("", searchHandler.Search)                      // 搜索资料
				search.GET("/hot-keywords", searchHandler.GetHotKeywords) // 热门搜索词
				search.GET("/history", searchHandler.GetSearchHistory)    // 搜索历史
				search.DELETE("/history", searchHandler.ClearSearchHistory) // 清空搜索历史
			}

			// 推荐相关
			recommendations := protected.Group("/materials")
			{
				recommendations.GET("/hot", searchHandler.GetHotMaterials)         // 热门资料
				recommendations.GET("/recommend", searchHandler.GetRecommendations) // 推荐资料
			}

			// 页面浏览记录（所有认证用户）
			protected.POST("/statistics/page-view", statisticsHandler.RecordPageView)
		}

		// 可选认证的路由（资料浏览，允许游客访问）
		optional := v1.Group("")
		optional.Use(middleware.OptionalJWTAuth(jwtManager, redisClient))
		{
			// 如果需要支持游客访问，可以取消以下注释
			// optional.GET("/materials", materialHandler.ListMaterials)          // 资料列表
			// optional.GET("/materials/search", materialHandler.SearchMaterials) // 搜索资料
			// optional.GET("/materials/:id", materialHandler.GetMaterial)        // 资料详情
		}
	}

	return r
}
