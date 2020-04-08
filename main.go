package konnek

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/functions/metadata"
	cloudevents "github.com/cloudevents/sdk-go"
	cloudeventsclient "github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	CloudEventsConsumer string `envconfig:"KONNEK_CE_CONSUMER" required:"true"`
}

type EventMetadata struct {
	ID     string
	Source string
	Type   string
}

func NewCloudEventsClient(cloudEventConsumer string) (cloudeventsclient.Client, error) {
	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(cloudEventConsumer),
		cloudeventshttp.WithEncoding(cloudeventshttp.Default),
	)
	if err != nil {
		return nil, err
	}

	client, err := cloudeventsclient.New(
		transport,
		cloudevents.WithDataContentType(cloudevents.ApplicationJSON),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Handler(ctx context.Context, event interface{}) error {
	log.Printf("context is: %+v", ctx)
	log.Printf("event is: %+v", event)

	eventMetadata, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("could not get metadata from context: %v", err)
	}

	var envConfig EnvConfig
	err = envconfig.Process("", &envConfig)
	if err != nil {
		return fmt.Errorf("could not load environment variables: %v", err)
	}

	// CE
	cloudEventsClient, err := NewCloudEventsClient(envConfig.CloudEventsConsumer)
	if err != nil {
		return fmt.Errorf("could not create client: %v", err)
	}

	cloudEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV1{
			ID:     eventMetadata.EventID,
			Source: *types.ParseURIRef(eventMetadata.Resource.Name),
			Type:   eventMetadata.EventType,
		},
		Data: event,
	}

	_, _, err = cloudEventsClient.Send(context.Background(), cloudEvent)
	if err != nil {
		return fmt.Errorf("could not send event: %v", err)
	}

	return nil
}
