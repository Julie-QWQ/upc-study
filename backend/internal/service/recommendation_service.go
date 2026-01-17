package service

import (
	"context"
	"fmt"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/repository"
	"gorm.io/gorm"
)

// recommendationService 推荐服务实现
type recommendationService struct {
	db             *gorm.DB
	materialRepo   repository.MaterialRepository
	downloadRepo   repository.DownloadRecordRepository
	favoriteRepo   repository.FavoriteRepository
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(
	db *gorm.DB,
	materialRepo repository.MaterialRepository,
	downloadRepo repository.DownloadRecordRepository,
	favoriteRepo repository.FavoriteRepository,
) RecommendationService {
	return &recommendationService{
		db:           db,
		materialRepo: materialRepo,
		downloadRepo: downloadRepo,
		favoriteRepo: favoriteRepo,
	}
}

// GetRecommendations 获取推荐资料
func (s *recommendationService) GetRecommendations(ctx context.Context, userID uint, req *model.RecommendationRequest) ([]*model.RecommendationResult, error) {
	switch req.Type {
	case "hot":
		// 热门资料推荐
		materials, err := s.GetHotMaterials(ctx, req.Limit)
		if err != nil {
			return nil, err
		}
		results := make([]*model.RecommendationResult, 0, len(materials))
		for _, m := range materials {
			results = append(results, &model.RecommendationResult{
				Material: m,
				Reason:   "热门资料",
				Score:    float64(m.DownloadCount+m.FavoriteCount*2) / 100.0,
			})
		}
		return results, nil

	case "personalized":
		// 个性化推荐
		return s.GetPersonalizedRecommendations(ctx, userID, req.Limit)

	case "related":
		// 相关资料推荐
		if req.MaterialID == nil {
			return nil, fmt.Errorf("相关推荐需要提供 material_id")
		}
		return s.GetRelatedMaterials(ctx, *req.MaterialID, req.Limit)

	case "downloaded":
		// 基于下载历史推荐
		return s.getRecommendedByDownloadHistory(ctx, userID, req.Limit)

	default:
		// 默认返回热门资料
		materials, err := s.GetHotMaterials(ctx, req.Limit)
		if err != nil {
			return nil, err
		}
		results := make([]*model.RecommendationResult, 0, len(materials))
		for _, m := range materials {
			results = append(results, &model.RecommendationResult{
				Material: m,
				Reason:   "热门资料",
				Score:    float64(m.DownloadCount+m.FavoriteCount*2) / 100.0,
			})
		}
		return results, nil
	}
}

// GetHotMaterials 获取热门资料
func (s *recommendationService) GetHotMaterials(ctx context.Context, limit int) ([]*model.Material, error) {
	var materials []*model.Material
	// 综合下载量、收藏量和浏览量计算热度
	err := s.db.WithContext(ctx).
		Where("status = ?", model.StatusApproved).
		Order("(download_count * 3 + favorite_count * 5 + view_count) DESC, created_at DESC").
		Limit(limit).
		Find(&materials).Error

	if err != nil {
		return nil, fmt.Errorf("获取热门资料失败: %w", err)
	}

	return materials, nil
}

// GetPersonalizedRecommendations 获取个性化推荐
func (s *recommendationService) GetPersonalizedRecommendations(ctx context.Context, userID uint, limit int) ([]*model.RecommendationResult, error) {
	// 1. 获取用户下载历史中的热门分类和课程
	downloads, _, err := s.downloadRepo.ListByUser(ctx, userID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取下载历史失败: %w", err)
	}

	// 2. 获取用户收藏
	favorites, _, err := s.favoriteRepo.ListByUser(ctx, userID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取收藏失败: %w", err)
	}

	// 3. 分析用户偏好
	categoryCount := make(map[string]int)
	courseCount := make(map[string]int)

	for _, d := range downloads {
		if d.Material != nil {
			categoryCount[string(d.Material.Category)]++
			if d.Material.CourseName != "" {
				courseCount[d.Material.CourseName]++
			}
		}
	}

	for _, f := range favorites {
		if f.Material != nil {
			categoryCount[string(f.Material.Category)]++
			if f.Material.CourseName != "" {
				courseCount[f.Material.CourseName]++
			}
		}
	}

	// 4. 找出最热门的分类和课程
	var topCategory string
	var topCourse string
	maxCount := 0

	for cat, count := range categoryCount {
		if count > maxCount {
			maxCount = count
			topCategory = cat
		}
	}

	maxCount = 0
	for course, count := range courseCount {
		if count > maxCount {
			maxCount = count
			topCourse = course
		}
	}

	// 5. 基于用户偏好推荐资料
	query := s.db.WithContext(ctx).
		Where("status = ?", model.StatusApproved).
		Order("download_count DESC, favorite_count DESC").
		Limit(limit)

	// 如果有偏好的分类，优先推荐该分类
	if topCategory != "" {
		query = query.Where("category = ?", topCategory)
	}
	// 如果有偏好的课程，也作为筛选条件
	if topCourse != "" {
		query = query.Where("course_name = ?", topCourse)
	}

	var materials []*model.Material
	if err := query.Find(&materials).Error; err != nil {
		return nil, fmt.Errorf("获取推荐资料失败: %w", err)
	}

	// 6. 构建推荐结果
	results := make([]*model.RecommendationResult, 0, len(materials))
	for _, m := range materials {
		reason := "根据您的浏览和下载历史推荐"
		if string(m.Category) == topCategory {
			reason = fmt.Sprintf("基于您喜欢的 %s 类资料推荐", m.Category)
		}
		if m.CourseName == topCourse {
			reason = fmt.Sprintf("基于 %s 课程相关资料推荐", m.CourseName)
		}

		results = append(results, &model.RecommendationResult{
			Material: m,
			Reason:   reason,
			Score:    0.8,
		})
	}

	return results, nil
}

// GetRelatedMaterials 获取相关资料
func (s *recommendationService) GetRelatedMaterials(ctx context.Context, materialID uint, limit int) ([]*model.RecommendationResult, error) {
	// 1. 获取原资料
	var material model.Material
	if err := s.db.WithContext(ctx).First(&material, materialID).Error; err != nil {
		return nil, fmt.Errorf("获取资料失败: %w", err)
	}

	// 2. 查找相关资料（相同分类或相同课程）
	var materials []*model.Material
	err := s.db.WithContext(ctx).
		Where("status = ? AND id != ?", model.StatusApproved, materialID).
		Where("(category = ? OR course_name = ?)", material.Category, material.CourseName).
		Order("download_count DESC, favorite_count DESC").
		Limit(limit).
		Find(&materials).Error

	if err != nil {
		return nil, fmt.Errorf("获取相关资料失败: %w", err)
	}

	// 3. 构建推荐结果
	results := make([]*model.RecommendationResult, 0, len(materials))
	for _, m := range materials {
		reason := "相关资料推荐"
		if m.Category == material.Category {
			reason = fmt.Sprintf("同为 %s 类资料", m.Category)
		}
		if m.CourseName == material.CourseName && m.CourseName != "" {
			reason = fmt.Sprintf("%s 课程相关资料", m.CourseName)
		}

		results = append(results, &model.RecommendationResult{
			Material: m,
			Reason:   reason,
			Score:    0.7,
		})
	}

	return results, nil
}

// getRecommendedByDownloadHistory 基于下载历史推荐
func (s *recommendationService) getRecommendedByDownloadHistory(ctx context.Context, userID uint, limit int) ([]*model.RecommendationResult, error) {
	// 1. 获取用户下载过的资料ID
	downloads, _, err := s.downloadRepo.ListByUser(ctx, userID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取下载历史失败: %w", err)
	}

	if len(downloads) == 0 {
		// 如果没有下载历史，返回热门资料
		materials, err := s.GetHotMaterials(ctx, limit)
		if err != nil {
			return nil, err
		}
		results := make([]*model.RecommendationResult, 0, len(materials))
		for _, m := range materials {
			results = append(results, &model.RecommendationResult{
				Material: m,
				Reason:   "热门资料",
				Score:    float64(m.DownloadCount+m.FavoriteCount*2) / 100.0,
			})
		}
		return results, nil
	}

	// 2. 提取已下载的资料ID
	downloadedIDs := make([]uint, len(downloads))
	for i, d := range downloads {
		downloadedIDs[i] = d.MaterialID
	}

	// 3. 查找其他下载了类似资料的用户也下载的资料（协同过滤）
	// 这里简化为：推荐与用户下载资料同分类或同课程的其他资料
	var materials []*model.Material
	err = s.db.WithContext(ctx).
		Raw(`
			SELECT DISTINCT m.* FROM materials m
			WHERE m.status = 'approved'
			AND m.id NOT IN (SELECT material_id FROM download_records WHERE user_id = ?)
			AND (
				m.category IN (SELECT DISTINCT category FROM materials WHERE id IN (SELECT material_id FROM download_records WHERE user_id = ?))
				OR m.course_name IN (SELECT DISTINCT course_name FROM materials WHERE id IN (SELECT material_id FROM download_records WHERE user_id = ?) AND course_name != '')
			)
			ORDER BY m.download_count DESC, m.favorite_count DESC
			LIMIT ?
		`, userID, userID, userID, limit).
		Scan(&materials).Error

	if err != nil {
		return nil, fmt.Errorf("获取推荐资料失败: %w", err)
	}

	// 4. 构建推荐结果
	results := make([]*model.RecommendationResult, 0, len(materials))
	for _, m := range materials {
		results = append(results, &model.RecommendationResult{
			Material: m,
			Reason:   "基于您的下载历史推荐",
			Score:    0.75,
		})
	}

	return results, nil
}
