package configuration

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

func BoolValueOr(ctx context.Context, key string, defaultValue bool) (bool, error) {
	log := logger.FromContext(ctx)
	log.Infof("getting bool value for key %s", key)
	colConfiguration, err := db.CollectionConfiguration(ctx)
	if err != nil {
		log.Errorf("failed to get configuration collection: %v", err)
		return false, err
	}
	res := colConfiguration.FindOne(ctx, bson.M{"key": key})
	configuration := new(db.Configuration)
	err = res.Decode(configuration)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Infof("key %s not found, returning default value %v", key, defaultValue)
			return defaultValue, nil
		}
		log.Errorf("failed to decode configuration: %v", err)
	}
	if configuration.BoolValue == nil {
		log.Infof("key %s found but value is nil, returning default value %v", key, defaultValue)
		return defaultValue, nil
	}
	log.Infof("key %s found, returning value %v", key, *configuration.BoolValue)
	return *configuration.BoolValue, nil
}
