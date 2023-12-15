package avmtester_test

import (
	"os"
	"path/filepath"
	"testing"

	terraform_module_test_helper "github.com/Azure/terraform-module-test-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestExamples(t *testing.T) {
	t.Parallel()
	modulePath := os.Getenv("AVM_MOD_PATH")
	if modulePath == "" {
		t.Skip("Cannot read AVM_MOD_PATH, you must set AVM_MOD_PATH to the avm module that you'd like to test. Skip examples test.")
	}
	files, err := os.ReadDir(filepath.Join(modulePath, "examples"))
	require.NoError(t, err)

	for _, file := range files {
		if file.IsDir() {
			exampleName := file.Name()
			t.Run(exampleName, func(t *testing.T) {
				t.Parallel()
				terraform_module_test_helper.RunE2ETest(t, modulePath, filepath.Join("examples", exampleName), terraform.Options{}, nil)
			})
		}
	}
}
