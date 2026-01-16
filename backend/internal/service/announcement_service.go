package service

import (
  "context"
  "errors"

  "github.com/study-upc/backend/internal/model"
  "github.com/study-upc/backend/internal/repository"
)

var (
  ErrAnnouncementNotFound = errors.New("announcement not found")
  ErrInvalidExpiresAt     = errors.New("expires_at must be after published_at")
)

// AnnouncementService provides announcement management.
type AnnouncementService interface {
  CreateAnnouncement(ctx context.Context, authorID uint, req *model.CreateAnnouncementRequest) (*model.AnnouncementResponse, error)
  GetAnnouncement(ctx context.Context, id uint) (*model.AnnouncementResponse, error)
  GetActiveAnnouncements(ctx context.Context, limit int) ([]model.Announcement, error)
  ListAnnouncements(ctx context.Context, req *model.AnnouncementListRequest) ([]model.Announcement, int64, error)
  UpdateAnnouncement(ctx context.Context, id, authorID uint, isAdmin bool, req *model.UpdateAnnouncementRequest) (*model.AnnouncementResponse, error)
  DeleteAnnouncement(ctx context.Context, id, authorID uint, isAdmin bool) error
}

type announcementService struct {
  announcementRepo repository.AnnouncementRepository
  userRepo         repository.UserRepository
}

func NewAnnouncementService(announcementRepo repository.AnnouncementRepository, userRepo repository.UserRepository) AnnouncementService {
  return &announcementService{
    announcementRepo: announcementRepo,
    userRepo:         userRepo,
  }
}

func (s *announcementService) CreateAnnouncement(ctx context.Context, authorID uint, req *model.CreateAnnouncementRequest) (*model.AnnouncementResponse, error) {
  if _, err := s.userRepo.FindByID(ctx, authorID); err != nil {
    return nil, err
  }

  if req.PublishedAt != nil && req.ExpiresAt != nil {
    if req.ExpiresAt.Before(*req.PublishedAt) || req.ExpiresAt.Equal(*req.PublishedAt) {
      return nil, ErrInvalidExpiresAt
    }
  }

  announcement := &model.Announcement{
    Title:       req.Title,
    Content:     req.Content,
    Priority:    req.Priority,
    AuthorID:    authorID,
    IsActive:    req.IsActive,
    PublishedAt: req.PublishedAt,
    ExpiresAt:   req.ExpiresAt,
  }

  if err := s.announcementRepo.Create(ctx, announcement); err != nil {
    return nil, err
  }

  announcement, err := s.announcementRepo.FindByID(ctx, announcement.ID)
  if err != nil {
    return nil, err
  }

  return announcement.ToAnnouncementResponse(), nil
}

func (s *announcementService) GetAnnouncement(ctx context.Context, id uint) (*model.AnnouncementResponse, error) {
  announcement, err := s.announcementRepo.FindByID(ctx, id)
  if err != nil {
    if errors.Is(err, repository.ErrAnnouncementNotFound) {
      return nil, ErrAnnouncementNotFound
    }
    return nil, err
  }

  return announcement.ToAnnouncementResponse(), nil
}

func (s *announcementService) GetActiveAnnouncements(ctx context.Context, limit int) ([]model.Announcement, error) {
  if limit <= 0 {
    limit = 5
  }
  if limit > 20 {
    limit = 20
  }

  return s.announcementRepo.FindActive(ctx, limit)
}

func (s *announcementService) ListAnnouncements(ctx context.Context, req *model.AnnouncementListRequest) ([]model.Announcement, int64, error) {
  return s.announcementRepo.List(ctx, req)
}

func (s *announcementService) UpdateAnnouncement(ctx context.Context, id, authorID uint, isAdmin bool, req *model.UpdateAnnouncementRequest) (*model.AnnouncementResponse, error) {
  announcement, err := s.announcementRepo.FindByID(ctx, id)
  if err != nil {
    if errors.Is(err, repository.ErrAnnouncementNotFound) {
      return nil, ErrAnnouncementNotFound
    }
    return nil, err
  }

  if announcement.AuthorID != authorID && !isAdmin {
    return nil, errors.New("not allowed to update this announcement")
  }

  if req.PublishedAt != nil && req.ExpiresAt != nil {
    if req.ExpiresAt.Before(*req.PublishedAt) || req.ExpiresAt.Equal(*req.PublishedAt) {
      return nil, ErrInvalidExpiresAt
    }
  }

  announcement.Title = req.Title
  announcement.Content = req.Content
  announcement.Priority = req.Priority
  announcement.IsActive = req.IsActive
  announcement.PublishedAt = req.PublishedAt
  announcement.ExpiresAt = req.ExpiresAt

  if err := s.announcementRepo.Update(ctx, announcement); err != nil {
    return nil, err
  }

  announcement, err = s.announcementRepo.FindByID(ctx, id)
  if err != nil {
    return nil, err
  }

  return announcement.ToAnnouncementResponse(), nil
}

func (s *announcementService) DeleteAnnouncement(ctx context.Context, id, authorID uint, isAdmin bool) error {
  announcement, err := s.announcementRepo.FindByID(ctx, id)
  if err != nil {
    if errors.Is(err, repository.ErrAnnouncementNotFound) {
      return ErrAnnouncementNotFound
    }
    return err
  }

  if announcement.AuthorID != authorID && !isAdmin {
    return errors.New("not allowed to delete this announcement")
  }

  return s.announcementRepo.Delete(ctx, id)
}
