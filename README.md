# AVM Tester

This Go module provides testing capabilities for [Azure Verified Modules (AVM)](https://aka.ms/avm) for Terraform. It is designed to be used from the Azure organisation via GitHub Actions.

## Purpose

The primary goals of `avmtester` are to:

1. **Test Module Examples (`avm_test.go`):** Verify that the examples provided within an AVM Terraform module apply cleanly without errors and are idempotent (meaning applying them multiple times results in the same state).
1. **Check for Breaking Changes (`avm_breaking_changes_test.go`):** Uses static code analysis to check for breaking changes by comparing the current version of a module against a specific prior Git tag (set via `PREVIOUS_TAG` environment variable).
1. **Test Destructive Upgrades (`avm_upgrade_destructive_test.go`):** Performs a deployment of the previous published version of the module and checks for destructive changes (ones involving a delete or replace operation).

## Integration

This tester is designed to be invoked via Make targets defined in the `avmmakefile` provided by the [Azure/tfmod-scaffold](https://github.com/Azure/tfmod-scaffold) project. The relevant targets are:

* `make test-example`: Runs the basic example tests (clean apply and idempotency).
* `make test-upgrade-destructive-example`: Runs the test that checks for destructive updates. (Assuming this is the correct target name based on the Go file)

These Make targets are typically called from GitHub Actions workflows.

When ran in CI, it uses the workflow inside the AVM template [.github/workflows/test-examples-template.yml](https://github.com/Azure/terraform-azurerm-avm-template/blob/main/.github/workflows/test-examples-template.yml).

This workflow sets up the necessary environment (including Azure credentials via OIDC) and runs the tests within a Docker container.

## Local Usage

While primarily intended for CI/CD, you can run these tests locally.

### Prerequisites

* Go installed
* Terraform installed
* Access to an Azure Subscription (Credentials configured via environment variables, OIDC, or other standard methods supported by the Azure provider)

### Environment Variables

Set the following environment variables before running the tests:

* `AVM_MOD_PATH`: **Required (All Tests).** Absolute or relative path to the root of the AVM module you want to test.
* `AVM_EXAMPLE`: **Required (Example & Destructive Upgrade Tests).** The name of the specific example directory within the module's `examples/` folder to test (e.g., `default`).
* `CURRENT_MAJOR_VERSION`: **Required (Destructive Upgrade Test).** The current major version number of the module (e.g., `0`, `1`). Used to find the previous version for destructive upgrade testing.
* `GITHUB_REPOSITORY`: **Required (Breaking Change & Destructive Upgrade Tests).** The GitHub repository where the module is hosted, in `owner/repo` format. Used to find previous versions/tags.
* `PREVIOUS_TAG`: **Required (Breaking Change Test).** The specific Git tag (e.g., `v1.0.0`) to compare the current code against for breaking changes.
* Azure Credentials: **Required (Example & Destructive Upgrade Tests).** Ensure your environment is configured for Terraform Azure provider authentication (e.g., `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_SUBSCRIPTION_ID`, `ARM_TENANT_ID`, or `ARM_USE_OIDC=true` with appropriate OIDC variables).

### Running Tests

1. Navigate to the `avmtester` directory in your terminal.
2. Source your environment variables (if saved in a file): `source .env`
3. Run the desired test:
    * **Example (successful apply & idempotency) Test:** `go test -v ./avm_test.go -timeout 30m`
    * **Breaking Change Test:** `go test -v ./avm_breaking_changes_test.go -timeout 10m`
    * **Destructive Upgrade Test:** `go test -v ./avm_upgrade_destructive_test.go -timeout 90m`

    *(Note: The `-timeout` flag increases the default test timeout, adjust as needed.)*

## Dependencies

This project uses [Azure/terraform-module-test-helper](https://github.com/Azure/terraform-module-test-helper): for testing logic, including the upgrade comparison.
