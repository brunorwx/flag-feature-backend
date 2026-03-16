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

func (r *InMemoryFeatureFlagRepository) GetAll() ([]*domain.FeatureFlag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flags := make([]*domain.FeatureFlag, 0, len(r.flags))
	for _, flag := range r.flags {
		flags = append(flags, flag)
	}
	return flags, nil
}
