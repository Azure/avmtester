package avmtester_test

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    terraform_module_test_helper "github.com/Azure/terraform-module-test-helper"
)

func TestBreakingChangeDetection(t *testing.T) {
    modulePath := os.Getenv("AVM_MOD_PATH")
    if modulePath == "" {
        t.Fatalf("Cannot read AVM_MOD_PATH, you must set AVM_MOD_PATH to the AVM module that you'd like to test.")
    }
    if _, err := os.Stat(modulePath); os.IsNotExist(err) {
        t.Fatalf("Module directory %s does not exist", modulePath)
    }

    repoEnv := os.Getenv("GITHUB_REPOSITORY")
    if repoEnv == "" {
        t.Fatalf("Error: GITHUB_REPOSITORY is not set. This must be set to the existing repository you want to test the upgrade against, using the 'org/repo' format.")
    }

    parts := strings.Split(repoEnv, "/")
    if len(parts) != 2 {
        t.Fatalf("Error: GITHUB_REPOSITORY '%s' is not in the expected 'org/repo' format", repoEnv)
    }
    githubOrg := parts[0]
    githubRepo := parts[1]
    if githubOrg == "" || githubRepo == "" {
        t.Fatalf("Error: GITHUB_REPOSITORY is in the expected 'org/repo' format, but either the org or repo is empty")
    } else {
        t.Logf("Running breaking change detection against GitHub Org: %s, Repo: %s\n", githubOrg, githubRepo)
    }

    previousTag := os.Getenv("PREVIOUS_TAG")
    if previousTag == "" {
        t.Fatalf("Cannot read PREVIOUS_TAG, you must set it to the Git tag you want to compare the current code against (e.g., 'v1.0.0').")
    }
    t.Logf("Comparing current module code against tag: %s", previousTag)

    // Run the breaking change detection test
    // Using filepath.Base for a reasonable sub-test name
    testName := filepath.Base(modulePath)
    t.Run(testName, func(t *testing.T) {
        changes, err := terraform_module_test_helper.BreakingChangesDetect(
            modulePath,
            githubOrg,
            githubRepo,
            &previousTag, // Pass the tag to compare against
        )

        if err != nil {
            t.Fatalf("Error during breaking change detection: %s", err)
        }

        if changes != "" {
            t.Fatalf("Breaking changes detected between current code and tag '%s':\n%s", previousTag, changes)
        }

        t.Logf("No breaking changes detected between current code and tag '%s'.", previousTag)
    })
}
