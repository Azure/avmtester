package avmtester_test

import (
	"os"
	"path/filepath"
	"testing"

	terraform_module_test_helper "github.com/Azure/terraform-module-test-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestExample(t *testing.T) {
	modulePath := os.Getenv("AVM_MOD_PATH")
	if modulePath == "" {
		t.Fatalf("Cannot read AVM_MOD_PATH, you must set AVM_MOD_PATH to the avm module that you'd like to test.")
	}
	example := os.Getenv("AVM_EXAMPLE")
	if modulePath == "" {
		t.Fatalf("Cannot read AVM_EXAMPLE, you must set AVM_EXAMPLE to the example name that you'd like to test.")
	}
	dir := filepath.Join(modulePath, "examples", example)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("Directory %s does not exist", dir)
	}
	terraform_module_test_helper.RunE2ETest(t, modulePath, filepath.Join("examples", example), terraform.Options{}, nil)
}
