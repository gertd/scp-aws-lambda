package kinesis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gertd/go-scp/scp"
	"github.com/gertd/scp-aws-lambda/pkg/config"
	"github.com/splunk/splunk-cloud-sdk-go/services/ingest"
)

// Handler --
func Handler(cfg *config.Config) func(ctx context.Context, kinesisEvent events.KinesisEvent) (err error) {

	clnt := scp.NewClient(cfg.Tenant, cfg.ClientID, cfg.ClientSecret)

	return func(ctx context.Context, kinesisEvent events.KinesisEvent) (err error) {

		for _, record := range kinesisEvent.Records {

			kinesisRecord := record.Kinesis
			data := kinesisRecord.Data

			var body interface{}

			if err := json.Unmarshal(data, &body); err != nil {
				return err
			}

			event := ingest.Event{
				Body:       body,
				Host:       &cfg.Host,
				Source:     &cfg.Source,
				Sourcetype: &cfg.SourceType,
			}

			if err := clnt.IngestEvent(&event); err != nil {
				return err
			}

			log.Printf("ingested kinesis seq# [%s]", record.Kinesis.SequenceNumber)
		}
		return nil
	}
}
