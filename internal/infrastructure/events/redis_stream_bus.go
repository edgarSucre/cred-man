package events

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/pkg/events"
	"github.com/redis/go-redis/v9"
)

type RedisStreamBus struct {
	client *redis.Client
	stream string
}

func (b *RedisStreamBus) Publish(ctx context.Context, event events.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.stream,
		Values: map[string]interface{}{
			"name": event.EventName(),
			"data": payload,
		},
	}).Err()
}
