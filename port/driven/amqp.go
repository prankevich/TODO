package driven

import "context"

type AmqpProducer interface {
	Publish(ctx context.Context, queueName string, msg []byte) error
}
