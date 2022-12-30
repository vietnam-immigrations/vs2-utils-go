package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/pkg/mongodb"
)

const colConfigName = "config"

func CollectionConfig(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/mongo/db", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colConfigName)
}

type Config struct {
	ID                              primitive.ObjectID `bson:"_id" json:"id"`
	PusherAppID                     string             `bson:"pusherAppId" json:"pusherAppId"`
	PusherAppKey                    string             `bson:"pusherAppKey" json:"pusherAppKey"`
	PusherAppSecret                 string             `bson:"pusherAppSecret" json:"pusherAppSecret"`
	PusherCluster                   string             `bson:"pusherCluster" json:"pusherCluster"`
	AdminDomain                     string             `bson:"adminDomain" json:"adminDomain"`
	CustomerDomain                  string             `bson:"customerDomain" json:"customerDomain"`
	EmailCustomerCC                 string             `bson:"emailCustomerCc" json:"emailCustomerCc"`
	EmailCustomerTemplateID         int                `bson:"emailCustomerTemplateId" json:"emailCustomerTemplateId"`
	EmailPartner                    string             `bson:"emailPartner" json:"emailPartner"`
	EmailPartnerCC                  string             `bson:"emailPartnerCc" json:"emailPartnerCc"`
	EmailPartnerTemplateID          int                `bson:"emailPartnerTemplateId" json:"emailPartnerTemplateId"`
	PriorityEmailCustomerTemplateID int                `bson:"priorityEmailCustomerTemplateId" json:"priorityEmailCustomerTemplateId"`
	PriorityEmailPartnerTemplateID  int                `bson:"priorityEmailPartnerTemplateId" json:"priorityEmailPartnerTemplateId"`
	APIRoot                         string             `bson:"apiRoot" json:"apiRoot"`
	APIKey                          string             `bson:"apiKey" json:"apiKey"`
	PushoverToken                   string             `bson:"pushoverToken" json:"pushoverToken"`
	PushoverUsers                   string             `bson:"pushoverUsers" json:"pushoverUsers"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	log := logger.FromContext(ctx)
	log.Infof("getting global configuration")
	colConfig, err := CollectionConfig(ctx)
	if err != nil {
		return nil, err
	}
	findAll, err := colConfig.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	items := make([]Config, 0)
	err = findAll.All(ctx, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		log.Warnf("missing global configuration, creating a new one")
		globalConfig := Config{
			ID: primitive.NewObjectID(),
		}
		_, err := colConfig.InsertOne(ctx, globalConfig)
		if err != nil {
			return nil, err
		}
		return &globalConfig, nil
	}
	if len(items) > 1 {
		return nil, fmt.Errorf("too many global configurations objects [%d]", len(items))
	}
	globalConfig := items[0]
	return &globalConfig, nil
}
