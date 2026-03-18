package infrastructure

import (
	"errors"
	"sync"

	"github.com/brunorwx/flagAPI/internal/domain"
)

type InMemoryFeatureFlagRepository struct {
	flags map[string]*domain.FeatureFlag
	mu    sync.RWMutex
}

func NewInMemoryFeatureFlagRepository() *InMemoryFeatureFlagRepository {
	return &InMemoryFeatureFlagRepository{
		flags: make(map[string]*domain.FeatureFlag),
	}
}

func (r *InMemoryFeatureFlagRepository) Save(flag *domain.FeatureFlag) error {
	if flag.Key == "" {
		return errors.New("feature flag key cannot be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flags[flag.Key] = flag
	return nil
}

func (r *InMemoryFeatureFlagRepository) GetByKey(key string) (*domain.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flag, ok := r.flags[key]
	if !ok {
		return nil, errors.New("feature flag not found")
	}
	return flag, nil
}

func (r *InMemoryFeatureFlagRepository) GetByKeys(keys []string) ([]*domain.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.FeatureFlag, 0)

	for _, key := range keys {
		if flag, ok := r.flags[key]; ok {
			result = append(result, flag)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no feature flags found")
	}

	return result, nil
}

func (r *InMemoryFeatureFlagRepository) GetAll() ([]*domain.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flags := make([]*domain.FeatureFlag, 0, len(r.flags))
	for _, flag := range r.flags {
		flags = append(flags, flag)
	}
	return flags, nil
}

func (r *InMemoryFeatureFlagRepository) DeleteByKey(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.flags[key]; !ok {
		return errors.New("feature flag not found")
	}
	delete(r.flags, key)
	return nil
}
