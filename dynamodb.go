package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)



var awsconfig aws.Config

func CreateDynamoClient() *dynamodb.Client {
	client := dynamodb.NewFromConfig(awsconfig)
	return dynamodb.New(sess)
}

func ConfigAWS() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	return cfg
}
