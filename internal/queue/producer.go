package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AuditProducer struct {
	channel *amqp.Channel
}

func NewAuditProducer(channel *amqp.Channel) *AuditProducer {
	return &AuditProducer{channel: channel}
}

func (p *AuditProducer) PublishLog(logData entity.AuditLog) error {
	_, err := p.channel.QueueDeclare(
		"audit_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctx,
		"",
		"audit_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	return err
}
