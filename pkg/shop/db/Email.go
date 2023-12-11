package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colShopEmail = "email"

func CollectionEmail(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/db-shop")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colShopEmail)
}

const MailingListAll = "all"

type Email struct {
	ID                       primitive.ObjectID `bson:"_id" json:"id"`
	Email                    string             `bson:"email" json:"email"`
	Secret                   string             `bson:"secret" json:"secret"`
	FirstName                string             `bson:"firstName" json:"firstName"`
	LastName                 string             `bson:"lastName" json:"lastName"`
	FullName                 string             `bson:"fullName" json:"fullName"`
	UnsubscribedMailingLists []string           `bson:"unsubscribedMailingLists" json:"unsubscribedMailingLists"`
	AlreadySentCampaigns     []string           `bson:"alreadySentCampaigns" json:"alreadySentCampaigns"`
	Complaints               []any              `bson:"complaints" json:"complaints"`
	CreatedAt                time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt                time.Time          `bson:"updatedAt" json:"updatedAt"`
}
