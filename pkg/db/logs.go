package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colLogsName = "logs"

func CollectionLogs(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/db")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colLogsName)
}

type LogType string

const (
	LogTypeOrder  = "order"
	LogTypeSystem = "system"
)

type Log struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Identifier string             `bson:"identifier" json:"identifier"`
	Type       LogType            `bson:"type" json:"type"`
	Message    string             `bson:"message" json:"message"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}
