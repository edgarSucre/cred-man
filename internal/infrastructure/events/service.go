package events

import (
	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNoRedisClient   = terror.Internal.New("redis-stream-bad-config", "redis client is missing")
	ErrInvalidStream   = terror.Internal.New("redis-stream-bad-config", "stream can't be empty")
	ErrInvalidGroup    = terror.Internal.New("redis-stream-bad-config", "consumer group can't be empty")
	ErrInvalidConsumer = terror.Internal.New("redis-stream-bad-config", "consumer can't be empty")
	ErrInvalidHandler  = terror.Internal.New("redis-stream-bad-config", "handlers can't be empty")
)

func NewStreamBus(redisClient *redis.Client, stream string) (*RedisStreamBus, error) {
	if redisClient == nil {
		return nil, ErrNoRedisClient
	}

	if len(stream) == 0 {
		return nil, ErrInvalidStream
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

func (params ConsumerParams) Validate() error {
	if params.Client == nil {
		return ErrNoRedisClient
	}

	if len(params.Group) == 0 {
		return ErrInvalidGroup
	}

	if len(params.Consumer) == 0 {
		return ErrInvalidConsumer
	}

	if len(params.Stream) == 0 {
		return ErrInvalidStream
	}

	if len(params.Handlers) == 0 {
		return ErrInvalidHandler
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
