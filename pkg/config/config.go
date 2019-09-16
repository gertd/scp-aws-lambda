package config

import (
	"fmt"
	"os"
)

const (
	// EnvInputType -- AWS input source type [SQS|SNS|KINESIS]
	EnvInputType = "SCP_INPUT_TYPE"
	// EnvTenant -- tenant name environment variable
	EnvTenant = "SCP_TENANT"
	// EnvClientID -- client id environment variable
	EnvClientID = "SCP_CLIENT_ID"
	// EnvClientSecret -- client secret environment variable
	EnvClientSecret = "SCP_CLIENT_SECRET"
	// EnvHost -- event host environment variable
	EnvHost = "SCP_EVENT_HOST"
	// EnvSource -- event source environment variable
	EnvSource = "SCP_EVENT_SOURCE"
	// EnvSourceType -- event source type environment variable
	EnvSourceType = "SCP_EVENT_SOURCETYPE"
)

const (
	validationErrorTemplate = "%s environment variable not set"
)

// Config --
type Config struct {
	EventSource  EventSourceType
	Tenant       string
	ClientID     string
	ClientSecret string
	Host         string
	Source       string
	SourceType   string
}

// NewConfigFromEnv --
func NewConfigFromEnv() *Config {
	return &Config{
		EventSource:  NewEventSourceType(os.Getenv(EnvInputType)),
		Tenant:       os.Getenv(EnvTenant),
		ClientID:     os.Getenv(EnvClientID),
		ClientSecret: os.Getenv(EnvClientSecret),
		Host:         os.Getenv(EnvHost),
		Source:       os.Getenv(EnvSource),
		SourceType:   os.Getenv(EnvSourceType),
	}
}

// Validate --
func (e *Config) Validate() error {

	if NewEventSourceType(os.Getenv(EnvInputType)) == EventSourceTypeUnknown {
		return fmt.Errorf(validationErrorTemplate, EnvInputType)
	}
	if os.Getenv(EnvClientID) == "" {
		return fmt.Errorf(validationErrorTemplate, EnvClientID)
	}
	if os.Getenv(EnvClientSecret) == "" {
		return fmt.Errorf(validationErrorTemplate, EnvClientSecret)
	}
	if os.Getenv(EnvTenant) == "" {
		return fmt.Errorf(validationErrorTemplate, EnvTenant)
	}
	return nil
}
