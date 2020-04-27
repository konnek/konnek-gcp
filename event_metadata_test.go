package konnek_gcp

import (
	"reflect"
	"testing"

	"cloud.google.com/go/functions/metadata"
	"golang.org/x/net/context"
)

func Test_getEventMetadata(t *testing.T) {
	r := metadata.Resource{
		Name: "event-source",
	}
	m := metadata.Metadata{
		EventID: "event-id",
		EventType: "event-type",
		Resource: &r,
	}
	parentCtx := context.Background()
	ctx := metadata.NewContext(parentCtx, &m)

	expected := &EventMetadata{
		Id: "event-id",
		Type: "event-type",
		Source: "event-source",
	}
	result, err := getEventMetadata(ctx)
	if err != nil {
		t.Errorf("expected err to be %v, got %v", nil, err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected eventMetadata to be %v, got %v", expected, result)
	}
}