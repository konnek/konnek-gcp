package konnek_gcp

import (
	"context"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Consumer string `envconfig:"KONNEK_CONSUMER" required:"true"`
}

func Handler(ctx context.Context, event interface{}) error {
	eventMetadata, err := getEventMetadata(ctx)
	if err != nil {
		return fmt.Errorf("could not get metadata from context: %v", err)
	}

	var env EnvConfig
	err = envconfig.Process("", &env)
	if err != nil {
		return fmt.Errorf("could not load environment variables: %v", err)
	}

	cloudEventsClient, err := newCloudEventsClient(env.Consumer)
	if err != nil {
		return fmt.Errorf("could not create client: %v", err)
	}

	cloudEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV1{
			ID:     eventMetadata.Id,
			Source: *types.ParseURIRef(eventMetadata.Source),
			Type:   eventMetadata.Type,
		},
		Data: event,
	}

	_, _, err = cloudEventsClient.Send(context.Background(), cloudEvent)
	if err != nil {
		return fmt.Errorf("could not send event: %v", err)
	}

	log.Printf("event with id %s send to %s", eventMetadata.Id, env.Consumer)

	return nil
}
