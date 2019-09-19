package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gertd/go-scp/ingest"
	"github.com/gertd/go-scp/scp"

	"github.com/gertd/scp-aws-lambda/cmd/scp-aws-lambda/kinesis"
	"github.com/gertd/scp-aws-lambda/cmd/scp-aws-lambda/sns"
	"github.com/gertd/scp-aws-lambda/cmd/scp-aws-lambda/sqs"
	"github.com/gertd/scp-aws-lambda/pkg/config"
)

// ldflags injected build version info
var (
	version string //nolint:gochecknoglobals
	date    string //nolint:gochecknoglobals
	commit  string //nolint:gochecknoglobals
)

func main() {
	log.Printf("%s %s-%s #%s %s ", version, runtime.GOOS, runtime.GOARCH, commit, date)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg := config.NewConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		return err
	}

	client := scp.NewClient(cfg.Tenant, cfg.ClientID, cfg.ClientSecret)
	if err := client.Authenticate(); err != nil {
		return err
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
		return fmt.Errorf("inputtype %s", cfg.EventSource)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// consume ingest events + produce batch evenys
	bp := ingest.NewBatchProcessor(ctx, cfg.Events)
	go bp.Run()

	// consume batch events
	bw := ingest.NewBatchWriter(ctx, client, bp.Batches())
	go bw.Run()

	// produce ingest events
	lambda.Start(handler)

	cancel()
	return nil
}
