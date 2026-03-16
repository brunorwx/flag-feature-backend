package domain

type FeatureFlag struct {
	Key           string
	Name          string
	GlobalEnabled bool
	UserOverrides map[string]bool
}

func NewFeatureFlag(key, name string, globalEnabled bool) *FeatureFlag {
	return &FeatureFlag{
		Key:           key,
		Name:          name,
		GlobalEnabled: globalEnabled,
		UserOverrides: make(map[string]bool),
	}
}

func (f *FeatureFlag) Evaluate(userID string) bool {
	if f.GlobalEnabled {
		return true
	} else if override, ok := f.UserOverrides[userID]; ok {
		return override
	}

	return false
}

func (f *FeatureFlag) SetUserOverride(userID string, enabled bool) {
	f.UserOverrides[userID] = enabled
}

func (f *FeatureFlag) SetGlobal(enabled bool) {
	f.GlobalEnabled = enabled
}
