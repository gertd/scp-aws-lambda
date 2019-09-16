package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gertd/scp-aws-lambda/cmd/scp-aws-lambda/kinesis"
	"github.com/gertd/scp-aws-lambda/cmd/scp-aws-lambda/sns"
	"github.com/gertd/scp-aws-lambda/cmd/scp-aws-lambda/sqs"
	"github.com/gertd/scp-aws-lambda/pkg/config"
)

func main() {
	cfg := config.NewConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	var handler interface{}

	switch cfg.EventSource {
	case config.EventSourceTypeKinesis:
		handler = kinesis.Handler(cfg)
	case config.EventSourceTypeSNS:
		handler = sns.Handler(cfg)
	case config.EventSourceTypeSQS:
		handler = sqs.Handler(cfg)
	default:
		log.Fatalf("inputtype %s", cfg.EventSource)
	}

	lambda.Start(handler)

}
