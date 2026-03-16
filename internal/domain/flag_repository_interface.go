package domain

type FeatureFlagRepository interface {
	Save(flag *FeatureFlag) error

	GetByKey(key string) (*FeatureFlag, error)

	GetAll() ([]*FeatureFlag, error)
}
