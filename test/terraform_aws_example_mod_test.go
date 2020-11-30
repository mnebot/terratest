package test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/stretchr/testify/assert"
)

var expectedName, expectedTagName string = "terratest-aws-example", "terratest-aws-example"
var awsRegion string = "us-east-1"
var expectedAMIName = "ubuntu-xenial-16.04-amd64-server"

/**
Test the EC2 deployment
*/
func TestEC2(t *testing.T) {
	t.Parallel()

	// execute Terraform
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-aws-example",
		Vars:         map[string]interface{}{},
		EnvVars:      map[string]string{"AWS_DEFAULT_REGION": awsRegion},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	instanceID := terraform.Output(t, terraformOptions, "instance_id")

	// Create EC2 client
	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if err != nil {
		fmt.Println(err.Error())
	}
	ec2svc := ec2.New(sess)
	ssmsvc := ssm.New(sess)

	// Run Checks
	CheckTags(t, ec2svc, instanceID)
	CheckAMI(t, ec2svc, instanceID)
	CheckSSMParameters(t, ssmsvc, instanceID)
}

// Check if expected tags exists
func CheckTags(t *testing.T, ec2svc *ec2.EC2, instanceID string) {
	// Look up the tags for the given Instance ID
	input := &ec2.DescribeTagsInput{Filters: []*ec2.Filter{{
		Name:   aws.String("resource-id"),
		Values: []*string{aws.String(instanceID)}}}}

	result, err := ec2svc.DescribeTags(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Check if the EC2 instance with a given tag and name is set.
	logger.Log(t, "Tag query result: "+string(result.String()))
	assert.Contains(t, result.String(), expectedTagName)
}

// Check if expected AMI is wich it is in the EC2
func CheckAMI(t *testing.T, ec2svc *ec2.EC2, instanceID string) {
	imageID := ""
	imageName := ""

	// Look up the AMI for the given Instance ID
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)}}

	result, err := ec2svc.DescribeInstances(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(result.Reservations) == 1 {
		imageID = *result.Reservations[0].Instances[0].ImageId
	}
	logger.Log(t, "ImageID: "+imageID)

	// Look up the AMI name for the given AMI ID
	imagesInput := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(imageID)}}

	resultImages, err := ec2svc.DescribeImages(imagesInput)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(resultImages.Images) == 1 {
		imageName = *resultImages.Images[0].Name
	}
	logger.Log(t, "Image Name: "+imageName)

	// Check if the EC2 image name is as expected.
	assert.Contains(t, imageName, expectedAMIName)
}

// Check if expected SSM parameter exists and have the correct value
func CheckSSMParameters(t *testing.T, ssmsvc ssmiface.SSMAPI, instanceID string) {

	name := "Test_EC2_instanceID"
	pname := &name

	results, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name: pname,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	logger.Log(t, "SSM Parameter Test_EC2_instanceID value: "+*results.Parameter.Value)

	// Check if the SSM Parameter with name "Test_EC2_instanceID" has the correct value
	assert.Equal(t, instanceID, *results.Parameter.Value)
}

// TODO S3 + KMS + naming (join) + TAGS corporativos