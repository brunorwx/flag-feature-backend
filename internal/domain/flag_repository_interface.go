package domain

type FeatureFlagRepository interface {
	Save(flag *FeatureFlag) error

	GetByKey(key string) (*FeatureFlag, error)

	GetByKeys(keys []string) ([]*FeatureFlag, error)

	GetAll() ([]*FeatureFlag, error)

	DeleteByKey(key string) error
}
