package avmtester_test

import (
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "testing"

    terraform_module_test_helper "github.com/Azure/terraform-module-test-helper"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestDestructiveUpgradeExample(t *testing.T) {
    // Read environment variables for module path and example name
    modulePath := os.Getenv("AVM_MOD_PATH")
    if modulePath == "" {
        t.Fatalf("Cannot read AVM_MOD_PATH, you must set AVM_MOD_PATH to the AVM module that you'd like to test.")
    }
    example := os.Getenv("AVM_EXAMPLE")
    if example == "" {
        t.Fatalf("Cannot read AVM_EXAMPLE, you must set AVM_EXAMPLE to the example name that you'd like to test.")
    }

    repoEnv := os.Getenv("GITHUB_REPOSITORY")
    if repoEnv == "" {
        t.Logf("Error: GITHUB_REPOSITORY is not set.  This must be set to the existing repository you want to test the upgrade against, using the 'org/repo' format.")
        os.Exit(1)
    }
    
    parts := strings.Split(repoEnv, "/")
    if len(parts) != 2 {
        t.Fatalf("Error: GITHUB_REPOSITORY '%s' is not in the expected 'org/repo' format", repoEnv)
        os.Exit(1)
    }   
    githubOrg := parts[0]
    githubRepo := parts[1]
    if githubOrg == "" || githubRepo == "" {
        t.Fatalf("Error: GITHUB_REPOSITORY is in the expected 'org/repo' format, but either the org or repo is empty")
    } else {
        t.Logf("Running upgrade against latest version published at GitHub Org: %s, Repo: %s\n", githubOrg, githubRepo)
    }

    // Construct the example directory path
    dir := filepath.Join(modulePath, "examples", example)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        t.Fatalf("Directory %s does not exist", dir)
    }

    // Check whether there are any *.tf files in the directory
    tfFiles, err := filepath.Glob(filepath.Join(dir, "*.tf"))
    if err != nil {
        t.Fatalf("Error while reading files in %s: %s", dir, err)
    }
    tfJsonFiles, err := filepath.Glob(filepath.Join(dir, "*.tf.json"))
    if err != nil {
        t.Fatalf("Error while reading files in %s: %s", dir, err)
    }
    if len(tfJsonFiles)+len(tfFiles) == 0 {
        t.Skipf("No Terraform files found in the directory, skip empty directory %s", example)
    }

    // Define Terraform options
    opts := terraform.Options{
        TerraformDir: dir,
    }

    // Get the current major version from the environment
    currentMajorVerStr := os.Getenv("CURRENT_MAJOR_VERSION")
    if currentMajorVerStr == "" {
        t.Fatalf("Cannot read CURRENT_MAJOR_VERSION, you must set it to the major version you want to test against (e.g., '1').")
    }
    currentMajorVer, err := strconv.Atoi(currentMajorVerStr)
    if err != nil {
        t.Fatalf("Invalid CURRENT_MAJOR_VERSION: %s", err)
    }

    // Run the upgrade test
    t.Run(example, func(t *testing.T) {
        terraform_module_test_helper.ModuleUpgradeDestructiveTest(
            t,
            githubOrg,
            githubRepo,
            filepath.Join("examples", example),
            modulePath,
            opts,
            currentMajorVer,
        )
    })
}