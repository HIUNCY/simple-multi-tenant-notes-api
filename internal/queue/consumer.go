package queue

import (
	"encoding/json"
	"log"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/entity"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AuditConsumer struct {
	channel   *amqp.Channel
	auditRepo repository.AuditRepository
}

func NewAuditConsumer(channel *amqp.Channel, repo repository.AuditRepository) *AuditConsumer {
	return &AuditConsumer{channel: channel, auditRepo: repo}
}

func (c *AuditConsumer) StartListening() {
	q, _ := c.channel.QueueDeclare("audit_queue", true, false, false, false, nil)

	msgs, err := c.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Gagal register consumer: %v", err)
	}

	go func() {
		log.Println("Worker mulai mendengarkan antrian 'audit_queue'...")
		for d := range msgs {
			var logData entity.AuditLog

			if err := json.Unmarshal(d.Body, &logData); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}

			err := c.auditRepo.Create(&logData)
			if err != nil {
				log.Printf("Gagal simpan ke Mongo: %v", err)
			} else {
				log.Printf("📥 [Worker] Sukses simpan log dari Queue: %s", logData.Action)
			}
		}
	}()
}
