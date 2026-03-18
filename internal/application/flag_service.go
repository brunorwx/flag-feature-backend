package application

import (
	"errors"

	"github.com/brunorwx/flagAPI/internal/domain"
)

type FeatureFlagService struct {
	repo domain.FeatureFlagRepository
}

func NewFeatureFlagService(repo domain.FeatureFlagRepository) *FeatureFlagService {
	return &FeatureFlagService{repo: repo}
}

func (s *FeatureFlagService) CreateFlag(key, name string, globalEnabled bool, rolloutPercentage int, targetRules [][]string) (*domain.FeatureFlag, error) {
	flag := domain.NewFeatureFlag(key, name, globalEnabled, rolloutPercentage, targetRules)
	if err := s.repo.Save(flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func (s *FeatureFlagService) GetFlag(key string) (*domain.FeatureFlag, error) {
	return s.repo.GetByKey(key)
}

func (s *FeatureFlagService) DeleteFlag(key string) error {
	if key == "" {
		return errors.New("flag key cannot be empty")
	}
	return s.repo.DeleteByKey(key)
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

func (s *FeatureFlagService) Evaluate(keys []string, userID string) ([]EvaluateResponse, error) {
	flags, err := s.repo.GetByKeys(keys)
	response := make([]EvaluateResponse, 0, len(flags))
	if err != nil {
		return response, err
	}

	for _, flag := range flags {
		response = append(response, EvaluateResponse{
			Key:     flag.Key,
			UserID:  userID,
			Enabled: flag.Evaluate(userID),
		})
	}
	return response, nil
}

func (s *FeatureFlagService) ListFlags() ([]*domain.FeatureFlag, error) {
	return s.repo.GetAll()
}
