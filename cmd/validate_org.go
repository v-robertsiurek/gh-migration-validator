/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"mona-actions/gh-migration-validator/internal/api"
	"mona-actions/gh-migration-validator/internal/validator"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// validateOrgCmd represents the validate-org command
var validateOrgCmd = &cobra.Command{
	Use:   "validate-org",
	Short: "Validate all repositories in an organization migration",
	Long: `Validate all repositories migrated from a source GitHub organization to a target
GitHub organization in a single run, producing one consolidated report.

This command:
- Lists all repositories in the source organization
- Validates each repository against the matching target repository
- Produces a single summary table with per-repo pass/fail/warn status
- Optionally writes a consolidated markdown report to a file

Repositories that cannot be accessed or fail validation are recorded with an
error but do not stop validation of the remaining repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate required variables
		if err := checkOrgVars(); err != nil {
			fmt.Printf("Organization validation configuration failed: %v\n", err)
			os.Exit(1)
		}

		sourceOrganization := viper.GetString("SOURCE_ORGANIZATION")
		targetOrganization := viper.GetString("TARGET_ORGANIZATION")

		// Initialize API with both source and target clients
		ghAPI, err := api.NewGitHubAPI()
		if err != nil {
			fmt.Printf("Failed to initialize API clients: %v\n", err)
			os.Exit(1)
		}

		// List repositories from the source organization
		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Listing repositories in %s...", sourceOrganization))
		repos, err := ghAPI.ListOrgRepos(api.SourceClient, sourceOrganization)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to list repositories: %v", err))
			os.Exit(1)
		}
		spinner.Success(fmt.Sprintf("Found %d repositories in %s", len(repos), sourceOrganization))

		if len(repos) == 0 {
			pterm.Warning.Println("No repositories found in source organization")
			return
		}

		// Run org-level validation
		migrationValidator := validator.New(ghAPI)
		summary := migrationValidator.ValidateOrganization(sourceOrganization, targetOrganization, repos)

		// Print consolidated results
		validator.PrintOrgValidationResults(summary)

		// Handle markdown output
		outputOrgMarkdown(summary)

		if viper.GetBool("STRICT_EXIT") && validator.OrgHasFailures(summary) {
			os.Exit(2)
		}
	},
}

func init() {
	// Source flags (local to this command)
	validateOrgCmd.Flags().StringP("github-source-org", "s", "", "Source GitHub organization")
	validateOrgCmd.Flags().StringP("github-source-pat", "a", "", "Source Organization GitHub token")
	validateOrgCmd.Flags().StringP("source-hostname", "u", "", "GitHub Enterprise source hostname url (optional)")

	// Bind to Viper in PreRun so they override env vars at execution time
	validateOrgCmd.PreRun = func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("SOURCE_ORGANIZATION", cmd.Flags().Lookup("github-source-org"))
		viper.BindPFlag("SOURCE_TOKEN", cmd.Flags().Lookup("github-source-pat"))
		viper.BindPFlag("SOURCE_HOSTNAME", cmd.Flags().Lookup("source-hostname"))
	}

	rootCmd.AddCommand(validateOrgCmd)
}

func checkOrgVars() error {
	required := map[string]requiredConfig{
		"SOURCE_ORGANIZATION": {"--github-source-org / -s", "GHMV_SOURCE_ORGANIZATION"},
		"TARGET_ORGANIZATION": {"--github-target-org / -t", "GHMV_TARGET_ORGANIZATION"},
		"SOURCE_TOKEN":        {"--github-source-pat / -a", "GHMV_SOURCE_TOKEN"},
		"TARGET_TOKEN":        {"--github-target-pat / -b", "GHMV_TARGET_TOKEN"},
	}

	for key, info := range required {
		if viper.GetString(key) == "" {
			return fmt.Errorf("%s is required. Set via %s flag or %s environment variable",
				key, info.flag, info.envVar)
		}
	}

	return nil
}

func outputOrgMarkdown(summary *validator.OrgValidationSummary) {
	markdownTable := viper.GetBool("MARKDOWN_TABLE")
	markdownFile := viper.GetString("MARKDOWN_FILE")

	if markdownTable {
		pterm.DefaultSection.Println("📋 Markdown Table (Copy-Paste Ready)")
		fmt.Println("```markdown")
		validator.WriteOrgMarkdownReport(summary, os.Stdout)
		fmt.Println("```")
		pterm.Info.Println("💡 Tip: You can select and copy the entire markdown section above to paste into documentation, issues, or pull requests!")
	}

	if markdownFile == "" {
		return
	}

	if dir := filepath.Dir(markdownFile); dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			pterm.Error.Printf("Directory %q does not exist for markdown file\n", dir)
			return
		}
	}

	var buf bytes.Buffer
	validator.WriteOrgMarkdownReport(summary, &buf)
	if err := os.WriteFile(markdownFile, buf.Bytes(), 0o644); err != nil {
		pterm.Error.Printf("Failed to write markdown file %s: %v\n", markdownFile, err)
		return
	}

	pterm.Success.Printf("📁 Markdown report saved to %s\n", markdownFile)
}
