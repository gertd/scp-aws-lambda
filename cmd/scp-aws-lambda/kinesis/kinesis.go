package kinesis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/gertd/scp-aws-lambda/pkg/config"
	"github.com/splunk/splunk-cloud-sdk-go/services/ingest"
)

// Handler --
func Handler(cfg *config.Config) func(ctx context.Context, kinesisEvent events.KinesisEvent) (err error) {

	return func(ctx context.Context, kinesisEvent events.KinesisEvent) (err error) {

		for _, record := range kinesisEvent.Records {

			kinesisRecord := record.Kinesis
			data := kinesisRecord.Data

			var body interface{}

			if err := json.Unmarshal(data, &body); err != nil {
				return err
			}

			ts := time.Now().UTC().Unix() * 1000
			ns := int32(0)

			event := ingest.Event{
				Body:       body,
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

			log.Printf("ingested kinesis seq# [%s]", record.Kinesis.SequenceNumber)
		}
		return nil
	}
}
