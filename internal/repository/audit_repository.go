package repository

import (
	"context"
	"time"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuditRepository interface {
	Create(log *entity.AuditLog) error
}

type mongoAuditRepository struct {
	client *mongo.Client
	dbName string
}

func NewMongoAuditRepository(client *mongo.Client, dbName string) AuditRepository {
	return &mongoAuditRepository{
		client: client,
		dbName: dbName,
	}
}

func (r *mongoAuditRepository) Create(log *entity.AuditLog) error {
	collection := r.client.Database(r.dbName).Collection("activity_logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, log)
	return err
}
