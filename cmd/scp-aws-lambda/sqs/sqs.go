package sqs

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gertd/scp-aws-lambda/pkg/config"
	"github.com/splunk/splunk-cloud-sdk-go/services/ingest"
)

// Handler --
func Handler(cfg *config.Config) func(ctx context.Context, sqsEvent events.SQSEvent) (err error) {

	return func(ctx context.Context, sqsEvent events.SQSEvent) (err error) {

		for _, message := range sqsEvent.Records {

			ts := time.Now().UTC().Unix() * 1000
			ns := int32(0)

			event := ingest.Event{
				Body:       &message.Body,
				Id:         &message.MessageId,
				Host:       &cfg.Host,
				Source:     &cfg.Source,
				Sourcetype: &cfg.SourceType,
				Timestamp:  &ts,
				Nanos:      &ns,
			}

			select {
			case cfg.Events <- event:
			case <-ctx.Done():
				return
			}

			log.Printf("ingested sqs msgid# [%s]", message.MessageId)
		}
		return nil
	}
}
