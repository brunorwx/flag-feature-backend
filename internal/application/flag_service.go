package application

import (
	"github.com/brunorwx/flagAPI/internal/domain"
)

type FeatureFlagService struct {
	repo domain.FeatureFlagRepository
}

func NewFeatureFlagService(repo domain.FeatureFlagRepository) *FeatureFlagService {
	return &FeatureFlagService{repo: repo}
}

func (s *FeatureFlagService) CreateFlag(key, name string, globalEnabled bool) (*domain.FeatureFlag, error) {
	flag := domain.NewFeatureFlag(key, name, globalEnabled)
	if err := s.repo.Save(flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func (s *FeatureFlagService) GetFlag(key string) (*domain.FeatureFlag, error) {
	return s.repo.GetByKey(key)
}

func (s *FeatureFlagService) SetUserOverride(key, userID string, enabled bool) (*domain.FeatureFlag, error) {
	flag, err := s.repo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	flag.SetUserOverride(userID, enabled)
	if err := s.repo.Save(flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func (s *FeatureFlagService) SetGlobalState(key string, enabled bool) (*domain.FeatureFlag, error) {
	flag, err := s.repo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	flag.SetGlobal(enabled)
	if err := s.repo.Save(flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func (s *FeatureFlagService) Evaluate(key, userID string) (bool, error) {
	flag, err := s.repo.GetByKey(key)
	if err != nil {
		return false, err
	}
	return flag.Evaluate(userID), nil
}

func (s *FeatureFlagService) ListFlags() ([]*domain.FeatureFlag, error) {
	return s.repo.GetAll()
}
