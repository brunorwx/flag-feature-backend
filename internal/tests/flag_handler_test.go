package tests

import (
	"testing"

	"github.com/brunorwx/flagAPI/internal/application"
	"github.com/brunorwx/flagAPI/internal/domain"
	"github.com/brunorwx/flagAPI/internal/infrastructure"
)

func TestDeleteFlag(t *testing.T) {
	repo := infrastructure.NewInMemoryFeatureFlagRepository()
	targetRules := [][]string{
		{"rule1", "rule2", "rule3"},
		{"ruleA", "ruleB"},
	}
	flag := domain.NewFeatureFlag("test-key", "Test", true, 50, targetRules)
	repo.Save(flag)

	service := application.NewFeatureFlagService(repo)

	err := service.DeleteFlag("test-key")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = service.GetFlag("test-key")
	if err == nil {
		t.Fatal("Expected error when getting deleted flag")
	}
}

func TestDeleteNonExistentFlag(t *testing.T) {
	repo := infrastructure.NewInMemoryFeatureFlagRepository()
	service := application.NewFeatureFlagService(repo)

	err := service.DeleteFlag("non-existent")
	if err == nil {
		t.Fatal("Expected error when deleting non-existent flag")
	}
}

func TestRolloutFlag(t *testing.T) {
	repo := infrastructure.NewInMemoryFeatureFlagRepository()

	targetRules := [][]string{
		{"rule1", "rule2", "rule3"},
		{"ruleA", "ruleB"},
	}

	flag := domain.NewFeatureFlag("test-key", "Test", false, 50, targetRules)
	flag2 := domain.NewFeatureFlag("test-key2", "Test", false, 50, targetRules)

	repo.Save(flag)
	repo.Save(flag2)

	service := application.NewFeatureFlagService(repo)

	responses, err := service.Evaluate([]string{"test-key", "test-key2"}, "user1")
	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}

	for _, res := range responses {
		if res.Enabled {
			t.Fatalf("expected user1 to NOT have access to flag %s", res.Key)
		}
	}

	responses2, err := service.Evaluate([]string{"test-key", "test-key2"}, "user2")
	if err != nil {
		t.Fatal(err)
	}

	for _, res := range responses2 {
		if !res.Enabled {
			t.Fatalf("expected user2 to have access to flag %s", res.Key)
		}
	}
}
