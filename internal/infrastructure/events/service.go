package events

import (
	"github.com/edgarSucre/mye"
	"github.com/redis/go-redis/v9"
)

//nolint:errcheck
func NewStreamBus(redisClient *redis.Client, stream string) (*RedisStreamBus, error) {
	err := mye.New(
		mye.CodeInternal,
		"stream_subscriber_config_error",
		"streamBus creation failed",
	)

	if redisClient == nil {
		err.WithField("redisClient", "redisClient is missing")
	}

	if len(stream) == 0 {
		err.WithField("stream", "stream name can't be empty")
	}

	if err.HasFields() {
		return nil, err
	}

	return &RedisStreamBus{client: redisClient, stream: stream}, nil
}

type ConsumerParams struct {
	Client   *redis.Client
	Consumer string
	Group    string
	Handlers map[string]EventHandler
	Stream   string
}

//nolint:errcheck
func (params ConsumerParams) Validate() error {
	err := mye.New(mye.CodeInternal, "redis_consumer_creation_failed", "parameter validation error")

	if params.Client == nil {
		err.WithField("client", "redis client is missing")
	}

	if len(params.Group) == 0 {
		err.WithField("group", "consumer group can't be empty")
	}

	if len(params.Consumer) == 0 {
		err.WithField("consumer", "consumer can't be empty")
	}

	if len(params.Stream) == 0 {
		err.WithField("stream", "stream can't be empty")
	}

	if len(params.Handlers) == 0 {
		err.WithField("handlers", "event handlers can't be empty")
	}

	if err.HasFields() {
		return err
	}

	return nil
}

func NewConsumer(params ConsumerParams) (*StreamConsumer, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &StreamConsumer{
		client:   params.Client,
		consumer: params.Consumer,
		group:    params.Group,
		handlers: params.Handlers,
		stream:   params.Stream,
		jobs:     make(chan redis.XMessage, 100),
	}, nil
}
