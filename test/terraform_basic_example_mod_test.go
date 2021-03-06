package test

import (
	"log"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the simple Terraform module in examples/terraform-basic-example using Terratest.
func TestTerraformBasicExample(t *testing.T) {
	t.Parallel()

	// Run `terraform output` to get the values of output variables
	actualTextExample := terraform.Output(t, terraformOptions, "example")
	actualTextExample2 := terraform.Output(t, terraformOptions, "example2")
	actualExampleList := terraform.OutputList(t, terraformOptions, "example_list")
	actualExampleMap := terraform.OutputMap(t, terraformOptions, "example_map")

	// website::tag::3::Check the output against expected values.
	// Verify we're getting back the outputs we expect
	assert.Equal(t, expectedText, actualTextExample)
	assert.Equal(t, expectedText, actualTextExample2)
	assert.Equal(t, expectedList, actualExampleList)
	assert.Equal(t, expectedMap, actualExampleMap)
}

func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	code := m.Run()
	log.Println("Do stuff AFTER the tests!")
	os.Exit(code)
}

// Init function called before tests have been executed
// Create AWS connection, initialize and apply Terraform
func Init() {
	expectedText := "test"
	expectedList := []string{expectedText}
	expectedMap := map[string]string{"expected": expectedText}

	terraformOptions := &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-basic-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"example": expectedText,

			// We also can see how lists and maps translate between terratest and terraform.
			"example_list": expectedList,
			"example_map":  expectedMap,
		},

		// Variables to pass to our Terraform code using -var-file options
		VarFiles: []string{"varfile.tfvars"},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	}

	// Run "terraform init" and "terraform apply".
	terraform.InitAndApply(t, terraformOptions)
}

// End function called afeter tests have been executed
// Clean up resources with "terraform destroy"
func End() {
	terraform.Destroy(t, terraformOptions)
}
