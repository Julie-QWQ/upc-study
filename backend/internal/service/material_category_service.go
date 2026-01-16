package service

import (
  "errors"

  "github.com/study-upc/backend/internal/model"
  "github.com/study-upc/backend/internal/repository"
  "gorm.io/gorm"
)

type MaterialCategoryService struct {
  categoryRepo *repository.MaterialCategoryRepository
}

func NewMaterialCategoryService(categoryRepo *repository.MaterialCategoryRepository) *MaterialCategoryService {
  return &MaterialCategoryService{
    categoryRepo: categoryRepo,
  }
}

func enrichCategory(category *model.MaterialCategoryConfig) *model.MaterialCategoryConfig {
  if category == nil {
    return nil
  }
  category.NameZh = category.Name
  category.NameEn = category.Code
  return category
}

func enrichCategoryList(categories []model.MaterialCategoryConfig) []model.MaterialCategoryConfig {
  for i := range categories {
    enrichCategory(&categories[i])
  }
  return categories
}

// List returns category list.
func (s *MaterialCategoryService) List(activeOnly bool) ([]model.MaterialCategoryConfig, error) {
  categories, err := s.categoryRepo.List(activeOnly)
  if err != nil {
    return nil, err
  }
  return enrichCategoryList(categories), nil
}

// GetByID returns a category by ID.
func (s *MaterialCategoryService) GetByID(id uint) (*model.MaterialCategoryConfig, error) {
  category, err := s.categoryRepo.GetByID(id)
  if err != nil {
    return nil, err
  }
  return enrichCategory(category), nil
}

// Create creates a new category.
func (s *MaterialCategoryService) Create(req model.MaterialCategoryRequest) (*model.MaterialCategoryConfig, error) {
  exists, err := s.categoryRepo.ExistsByCode(req.Code, 0)
  if err != nil {
    return nil, err
  }
  if exists {
    return nil, errors.New("category code already exists")
  }

  category := &model.MaterialCategoryConfig{
    Code:        req.Code,
    Name:        req.Name,
    Description: req.Description,
    Icon:        req.Icon,
    SortOrder:   req.SortOrder,
  }

  if req.IsActive != nil {
    category.IsActive = *req.IsActive
  } else {
    category.IsActive = true
  }

  if err := s.categoryRepo.Create(category); err != nil {
    return nil, err
  }

  return enrichCategory(category), nil
}

// Update updates a category.
func (s *MaterialCategoryService) Update(id uint, req model.MaterialCategoryRequest) (*model.MaterialCategoryConfig, error) {
  category, err := s.categoryRepo.GetByID(id)
  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return nil, errors.New("category not found")
    }
    return nil, err
  }

  exists, err := s.categoryRepo.ExistsByCode(req.Code, id)
  if err != nil {
    return nil, err
  }
  if exists {
    return nil, errors.New("category code already exists")
  }

  category.Code = req.Code
  category.Name = req.Name
  category.Description = req.Description
  category.Icon = req.Icon
  category.SortOrder = req.SortOrder

  if req.IsActive != nil {
    category.IsActive = *req.IsActive
  }

  if err := s.categoryRepo.Update(category); err != nil {
    return nil, err
  }

  return enrichCategory(category), nil
}

// Delete deletes a category.
func (s *MaterialCategoryService) Delete(id uint) error {
  category, err := s.categoryRepo.GetByID(id)
  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return errors.New("category not found")
    }
    return err
  }

  count, err := s.categoryRepo.CheckUsage(category.Code)
  if err != nil {
    return err
  }
  if count > 0 {
    return errors.New("category is in use and cannot be deleted")
  }

  return s.categoryRepo.Delete(id)
}

// ToggleStatus toggles category active status.
func (s *MaterialCategoryService) ToggleStatus(id uint) (*model.MaterialCategoryConfig, error) {
  category, err := s.categoryRepo.GetByID(id)
  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return nil, errors.New("category not found")
    }
    return nil, err
  }

  category.IsActive = !category.IsActive
  if err := s.categoryRepo.Update(category); err != nil {
    return nil, err
  }

  return enrichCategory(category), nil
}
