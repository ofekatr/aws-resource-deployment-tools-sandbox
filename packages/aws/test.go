package aws

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func Test() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	getResources(cfg)

	createResource(cfg)

	updateResource(cfg)

	deleteResource(cfg)
}

func getResources(cfg aws.Config) {
	client := resourcegroupstaggingapi.NewFromConfig(cfg)

	output, err := client.GetResources(context.Background(), &resourcegroupstaggingapi.GetResourcesInput{})

	if err != nil {
		log.Fatalf("failed to get resources, %v", err)
	}

	jsonOutput, err := json.Marshal(*output)

	log.Printf("Get resources output: %v\n", string(jsonOutput))
}

func createResource(cfg aws.Config) {
	path := filepath.Join("packages", "aws", "fixtures", "cloudformation", "s3.yaml")

	templateBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read template file, %v", err)
	}

	log.Println("Template: ", string(templateBytes))

	mySession := session.Must(session.NewSession())

	// Create a CloudFormation client from just a session.
	svc := cloudformation.New(mySession)

	output, err := svc.CreateStack(&cloudformation.CreateStackInput{ // TODO: How should this be polled? When should we scan for issues?
		StackName:    aws.String("test-kms-stack"),
		TemplateBody: aws.String(string(templateBytes)), // TODO: See if you pass the file path instead
	})
	if err != nil {
		log.Fatalf("failed to create s3 bucket, %v", err)
	}

	log.Printf("Create s3 bucket output: %v\n", output)

	if svc.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String("test-kms-stack"),
	}); err != nil {
		log.Fatalf("failed to wait for stack to be created, %v", err)
	}
}

func updateResource(cfg aws.Config) {
	path := filepath.Join("packages", "aws", "fixtures", "cloudformation", "s3-ignore-public-acls.yaml")

	templateBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read template file, %v", err)
	}

	log.Println("Template: ", string(templateBytes))

	mySession := session.Must(session.NewSession())

	// Create a CloudFormation client from just a session.
	svc := cloudformation.New(mySession)

	output, err := svc.UpdateStack(&cloudformation.UpdateStackInput{
		StackName:    aws.String("test-kms-stack"),
		TemplateBody: aws.String(string(templateBytes)),
	})
	if err != nil {
		log.Fatalf("failed to update s3 bucket, %v", err)
	}

	log.Printf("Update s3 output: %v\n", output)

	if svc.WaitUntilStackUpdateComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String("test-kms-stack"),
	}); err != nil {
		log.Fatalf("failed to wait for stack to be updated, %v", err)
	}
}

func deleteResource(cfg aws.Config) {
	mySession := session.Must(session.NewSession())

	// Create a CloudFormation client from just a session.
	svc := cloudformation.New(mySession)

	output, err := svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String("test-kms-stack"),
	})
	if err != nil {
		log.Fatalf("failed to delete s3 bucket, %v", err)
	}

	log.Printf("Delete s3 output: %v\n", output)
}
