package sns

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gertd/scp-aws-lambda/pkg/config"
	"github.com/splunk/splunk-cloud-sdk-go/services/ingest"
)

// Handler --
func Handler(cfg *config.Config) func(ctx context.Context, snsEvent events.SNSEvent) (err error) {

	return func(ctx context.Context, snsEvent events.SNSEvent) (err error) {

		for _, record := range snsEvent.Records {

			snsRecord := record.SNS

			ts := time.Now().UTC().Unix() * 1000
			ns := int32(0)

			event := ingest.Event{
				Body:       &snsRecord.Message,
				Id:         &snsRecord.MessageID,
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

			log.Printf("ingested sns msgid# [%s]", snsRecord.MessageID)
		}
		return nil
	}
}
