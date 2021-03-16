package pubsub

import (
	"testing"

	"github.com/matryer/is"
)

func TestDoublePublish(t *testing.T) {
	is := is.New(t)

	_, err := Publish(AUTO_BACKUP)
	is.NoErr(err) // can publish topic
	_, err = Publish(AUTO_BACKUP)
	is.True(err != nil) // cant  re-publish same topic
}
