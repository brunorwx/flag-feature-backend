package domain

import (
	"hash/fnv"
	"time"
)

type FeatureFlag struct {
	Key               string
	Name              string
	GlobalEnabled     bool
	UserOverrides     map[string]bool
	RolloutPercentage int
	Timestamp         time.Time
	TargetRules       TargetRule
}

func NewFeatureFlag(key, name string, globalEnabled bool, rolloutPercentage int, targetRules [][]string) *FeatureFlag {
	return &FeatureFlag{
		Key:               key,
		Name:              name,
		GlobalEnabled:     globalEnabled,
		UserOverrides:     make(map[string]bool),
		RolloutPercentage: rolloutPercentage,
		Timestamp:         time.Now(),
	}
}

func (f *FeatureFlag) Evaluate(userID string) bool {
	if override, ok := f.UserOverrides[userID]; ok {
		return override
	}
	if f.GlobalEnabled {
		return true
	}
	if f.RolloutPercentage == 0 {
		return false
	}
	if f.RolloutPercentage >= 100 {
		return true
	}

	userBucket := hashUserID(userID) % 100
	isInRollout := userBucket < f.RolloutPercentage
	return isInRollout

}

func hashUserID(userID string) int {
	hashedID := fnv.New32a()
	hashedID.Write([]byte(userID))
	return int(hashedID.Sum32())
}

func (f *FeatureFlag) SetUserOverride(userID string, enabled bool) {
	f.UserOverrides[userID] = enabled
}

func (f *FeatureFlag) SetGlobal(enabled bool) {
	f.GlobalEnabled = enabled
}
