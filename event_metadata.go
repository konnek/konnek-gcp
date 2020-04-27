package konnek_gcp

import (
	"context"

	"cloud.google.com/go/functions/metadata"
)

type EventMetadata struct {
	Type   string
	Source string
	Id     string
}

func getEventMetadata(ctx context.Context) (*EventMetadata, error) {
	md, err := metadata.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	eventMetadata := EventMetadata{
		Type:   md.EventType,
		Id:     md.EventID,
		Source: md.Resource.Name,
	}

	return &eventMetadata, nil
}
