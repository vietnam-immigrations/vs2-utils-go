package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colShopConfiguration = "configuration"

func CollectionConfiguration(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/db-shop")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colShopConfiguration)
}

type Configuration struct {
	ID        string `bson:"_id" json:"id"`
	Key       string `bson:"key" json:"key"`
	BoolValue *bool  `bson:"boolValue,omitempty" json:"boolValue,omitempty"`
}
