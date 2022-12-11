package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/mongodb"
)

const colResultsName = "results"

func CollectionResult(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/mongo/db", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colResultsName)
}

type ResultFile struct {
	Name         string `bson:"name" json:"name"`
	Processed    bool   `bson:"processed" json:"processed"`
	ErrorMessage string `bson:"errorMessage" json:"errorMessage"`
	OrderNumber  string `bson:"orderNumber" json:"orderNumber"`
	// PassportNumber used to match CV manually
	PassportNumber string `bson:"passportNumber" json:"passportNumber"`
}

type Result struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	S3DirKey string             `bson:"s3DirKey" json:"s3DirKey"`
	Files    []ResultFile       `bson:"files" json:"files"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
