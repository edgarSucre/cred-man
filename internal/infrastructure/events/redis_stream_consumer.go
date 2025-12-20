package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	StreamConsumer struct {
		client   *redis.Client
		group    string
		consumer string
		handlers map[string]EventHandler
	}

	EventEnvelope struct {
		ID         string          `json:"id"`
		Name       string          `json:"name"`
		Payload    json.RawMessage `json:"payload"`
		OccurredAt time.Time       `json:"occurred_at"`
	}

	EventHandler interface {
		EventName() string
		Handle(ctx context.Context, payload json.RawMessage) error
	}
)

const streamName = "domain-events"

func (c *StreamConsumer) Start(ctx context.Context) {
	for {
		streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.group,
			Consumer: c.consumer,
			Streams:  []string{streamName, ">"},
			Block:    time.Second,
			Count:    10,
		}).Result()

		if err != nil {
			continue
		}

		for _, msg := range streams[0].Messages {
			c.processMessage(ctx, msg)
		}
	}
}

func (c *StreamConsumer) processMessage(
	ctx context.Context,
	msg redis.XMessage,
) {
	var envelope EventEnvelope
	raw, ok := msg.Values["data"].(string)
	if !ok {
		// malformed message → ACK to avoid poison loop
		c.ack(ctx, msg)
		return
	}

	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		// invalid JSON → ACK or move to DLQ
		c.ack(ctx, msg)
		return
	}

	handler, ok := c.handlers[envelope.Name]
	if !ok {
		// unknown event → ack & ignore
		c.ack(ctx, msg)
		return
	}

	if err := handler.Handle(ctx, envelope.Payload); err != nil {
		// do NOT ack → retry later
		return
	}

	c.ack(ctx, msg)
}

func (c *StreamConsumer) ack(ctx context.Context, msg redis.XMessage) {
	_ = c.client.XAck(
		ctx,
		streamName,
		c.group,
		msg.ID,
	).Err()
}
