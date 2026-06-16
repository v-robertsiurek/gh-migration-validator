package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestCheckOrgVars_AllProvided(t *testing.T) {
	resetViperAndEnv()
	defer resetViperAndEnv()

	os.Setenv("GHMV_SOURCE_ORGANIZATION", "source-org")
	os.Setenv("GHMV_TARGET_ORGANIZATION", "target-org")
	os.Setenv("GHMV_SOURCE_TOKEN", "source-token")
	os.Setenv("GHMV_TARGET_TOKEN", "target-token")

	viper.SetEnvPrefix("GHMV")
	viper.AutomaticEnv()

	err := checkOrgVars()
	if err != nil {
		t.Errorf("Expected no error when all vars are provided, got: %v", err)
	}
}

func TestCheckOrgVars_MissingSourceOrg(t *testing.T) {
	resetViperAndEnv()
	defer resetViperAndEnv()

	os.Setenv("GHMV_TARGET_ORGANIZATION", "target-org")
	os.Setenv("GHMV_SOURCE_TOKEN", "token")
	os.Setenv("GHMV_TARGET_TOKEN", "token")

	viper.SetEnvPrefix("GHMV")
	viper.AutomaticEnv()

	err := checkOrgVars()
	if err == nil {
		t.Error("Expected error when SOURCE_ORGANIZATION is missing")
	}
	if !strings.Contains(err.Error(), "SOURCE_ORGANIZATION") {
		t.Errorf("Error should mention SOURCE_ORGANIZATION, got: %v", err)
	}
}

func TestCheckOrgVars_MissingTargetOrg(t *testing.T) {
	resetViperAndEnv()
	defer resetViperAndEnv()

	os.Setenv("GHMV_SOURCE_ORGANIZATION", "source-org")
	os.Setenv("GHMV_SOURCE_TOKEN", "token")
	os.Setenv("GHMV_TARGET_TOKEN", "token")

	viper.SetEnvPrefix("GHMV")
	viper.AutomaticEnv()

	err := checkOrgVars()
	if err == nil {
		t.Error("Expected error when TARGET_ORGANIZATION is missing")
	}
	if !strings.Contains(err.Error(), "TARGET_ORGANIZATION") {
		t.Errorf("Error should mention TARGET_ORGANIZATION, got: %v", err)
	}
}

func TestCheckOrgVars_MissingSourceToken(t *testing.T) {
	resetViperAndEnv()
	defer resetViperAndEnv()

	os.Setenv("GHMV_SOURCE_ORGANIZATION", "source-org")
	os.Setenv("GHMV_TARGET_ORGANIZATION", "target-org")
	os.Setenv("GHMV_TARGET_TOKEN", "token")

	viper.SetEnvPrefix("GHMV")
	viper.AutomaticEnv()

	err := checkOrgVars()
	if err == nil {
		t.Error("Expected error when SOURCE_TOKEN is missing")
	}
	if !strings.Contains(err.Error(), "SOURCE_TOKEN") {
		t.Errorf("Error should mention SOURCE_TOKEN, got: %v", err)
	}
}

func TestCheckOrgVars_MissingTargetToken(t *testing.T) {
	resetViperAndEnv()
	defer resetViperAndEnv()

	os.Setenv("GHMV_SOURCE_ORGANIZATION", "source-org")
	os.Setenv("GHMV_TARGET_ORGANIZATION", "target-org")
	os.Setenv("GHMV_SOURCE_TOKEN", "token")

	viper.SetEnvPrefix("GHMV")
	viper.AutomaticEnv()

	err := checkOrgVars()
	if err == nil {
		t.Error("Expected error when TARGET_TOKEN is missing")
	}
	if !strings.Contains(err.Error(), "TARGET_TOKEN") {
		t.Errorf("Error should mention TARGET_TOKEN, got: %v", err)
	}
}

func TestCheckOrgVars_DoesNotRequireRepoFlags(t *testing.T) {
	resetViperAndEnv()
	defer resetViperAndEnv()

	// Org validation should NOT require SOURCE_REPO / TARGET_REPO
	os.Setenv("GHMV_SOURCE_ORGANIZATION", "source-org")
	os.Setenv("GHMV_TARGET_ORGANIZATION", "target-org")
	os.Setenv("GHMV_SOURCE_TOKEN", "token")
	os.Setenv("GHMV_TARGET_TOKEN", "token")

	viper.SetEnvPrefix("GHMV")
	viper.AutomaticEnv()

	err := checkOrgVars()
	if err != nil {
		t.Errorf("checkOrgVars should not require repo flags, got: %v", err)
	}
}

func TestValidateOrgCommand_Registered(t *testing.T) {
	// Verify the validate-org subcommand is registered on rootCmd
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "validate-org" {
			found = true
			break
		}
	}
	if !found {
		t.Error("validate-org command should be registered as a subcommand of root")
	}
}
