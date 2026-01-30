package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	StreamConsumer struct {
		client   *redis.Client
		consumer string
		group    string
		handlers map[string]EventHandler
		jobs     chan redis.XMessage
		stream   string
	}

	EventEnvelope struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		Payload    string    `json:"data"`
		OccurredAt time.Time `json:"occurred_at"`
	}

	EventHandler interface {
		EventName() string
		Handle(ctx context.Context, payload json.RawMessage) error
	}
)

// simple workers implementation
func (c *StreamConsumer) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.jobs:
			if !ok {
				return // channel closed
			}
			c.processMessage(ctx, msg)
		}
	}
}

func (c *StreamConsumer) Start(ctx context.Context) error {
	// start workers
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			c.worker(ctx)
		}()
	}

	for {

		select {
		case <-ctx.Done():
			// stop accepting new work
			close(c.jobs)

			// wait for workers to finish
			wg.Wait()
			return nil
		default:
			streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    c.group,
				Consumer: c.consumer,
				Streams:  []string{c.stream, ">"},
				Block:    time.Second,
				Count:    20,
			}).Result()

			if err != nil {
				continue
			}

			for _, msg := range streams[0].Messages {
				// implicit backpressure when there are no workers available
				select {
				case c.jobs <- msg:
				case <-ctx.Done():
					// stop accepting new work
					close(c.jobs)

					// wait for workers to finish
					wg.Wait()
					return nil
				}
			}
		}
	}
}

func (c *StreamConsumer) processMessage(
	ctx context.Context,
	msg redis.XMessage,
) {
	var envelope EventEnvelope

	raw, err := json.Marshal(msg.Values)
	if err != nil {
		// malformed message, ACK to remove it from stream
		c.ack(ctx, msg)
		return
	}

	if err := json.Unmarshal(raw, &envelope); err != nil {
		// invalid JSON, nothing to do here, ACK to remove it from stream
		c.ack(ctx, msg)
		return
	}

	var payload json.RawMessage

	if err := json.Unmarshal([]byte(envelope.Payload), &payload); err != nil {
		// invalid JSON, nothing to do here, ACK to remove it from stream
		c.ack(ctx, msg)
		return
	}

	handler, ok := c.handlers[envelope.Name]
	if !ok {
		// unknown event, ACK to remove it from stream
		c.ack(ctx, msg)
		return
	}

	if err := handler.Handle(ctx, payload); err != nil {
		// handler error, try again
		return
	}

	c.ack(ctx, msg)
}

func (c *StreamConsumer) ack(ctx context.Context, msg redis.XMessage) {
	_ = c.client.XAck(
		ctx,
		c.stream,
		c.group,
		msg.ID,
	).Err()
}
