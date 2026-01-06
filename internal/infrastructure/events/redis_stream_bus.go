package events

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/edgarSucre/crm/internal/domain/event"
	"github.com/edgarSucre/mye"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisStreamBus struct {
	client *redis.Client
	stream string
}

const eventPublishErrSlug = "event_publishing_failed"

func (b *RedisStreamBus) Publish(ctx context.Context, event event.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return mye.Wrap(err, mye.CodeInvalid, eventPublishErrSlug, "json marshal error").
			WithAttribute("event", event)
	}

	err = b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.stream,
		Values: map[string]interface{}{
			"id":          uuid.New().String(),
			"name":        event.EventName(),
			"data":        payload,
			"occurred_at": time.Now(),
		},
	}).Err()

	if err != nil {
		if errIsTemporary(err) {
			return mye.Wrap(err, mye.CodeTimeout, eventPublishErrSlug, "redis connection error").
				WithUserMsg("We are experiencing delay issues due to high traffic. Please try again in a few seconds")
		}

		return mye.Wrap(err, mye.CodeInternal, eventPublishErrSlug, "redis internal error")
	}

	return err
}

func errIsTemporary(err error) bool {
	return errors.Is(err, redis.ErrClosed) ||
		errors.Is(err, redis.ErrPoolExhausted) ||
		errors.Is(err, redis.ErrPoolTimeout) ||
		errors.Is(err, context.DeadlineExceeded)
}
