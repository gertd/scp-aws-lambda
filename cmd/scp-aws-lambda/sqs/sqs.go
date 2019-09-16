package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gertd/go-scp/scp"
	"github.com/gertd/scp-aws-lambda/pkg/config"
	"github.com/splunk/splunk-cloud-sdk-go/services/ingest"
)

// Handler --
func Handler(cfg *config.Config) func(ctx context.Context, sqsEvent events.SQSEvent) (err error) {

	clnt := scp.NewClient(cfg.Tenant, cfg.ClientID, cfg.ClientSecret)

	return func(ctx context.Context, sqsEvent events.SQSEvent) (err error) {

		for _, message := range sqsEvent.Records {

			event := ingest.Event{
				Body:       &message.Body,
				Id:         &message.MessageId,
				Host:       &cfg.Host,
				Source:     &cfg.Source,
				Sourcetype: &cfg.SourceType,
			}

			if err := clnt.IngestEvent(&event); err != nil {
				return err
			}

			log.Printf("ingested sqs msgid# [%s]", message.MessageId)
		}
		return nil
	}
}
