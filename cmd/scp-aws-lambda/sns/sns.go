package sns

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gertd/go-scp/scp"
	"github.com/gertd/scp-aws-lambda/pkg/config"
	"github.com/splunk/splunk-cloud-sdk-go/services/ingest"
)

// Handler --
func Handler(cfg *config.Config) func(ctx context.Context, snsEvent events.SNSEvent) (err error) {

	clnt := scp.NewClient(cfg.Tenant, cfg.ClientID, cfg.ClientSecret)

	return func(ctx context.Context, snsEvent events.SNSEvent) (err error) {

		for _, record := range snsEvent.Records {

			snsRecord := record.SNS

			event := ingest.Event{
				Body:       &snsRecord.Message,
				Id:         &snsRecord.MessageID,
				Host:       &cfg.Host,
				Source:     &cfg.Source,
				Sourcetype: &cfg.SourceType,
			}

			if err := clnt.IngestEvent(&event); err != nil {
				return err
			}

			log.Printf("ingested sns msgid# [%s]", snsRecord.MessageID)
		}
		return nil
	}
}
