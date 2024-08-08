package logs

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

func Add(ctx context.Context, logType db.LogType, identifier string, message string) {
	log := logger.FromContext(ctx)
	log.Infof("Add log %s [%s]: %s", identifier, logType, message)

	colLogs, err := db.CollectionLogs(ctx)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	colLogs.InsertOne(ctx, db.Log{
		ID:         primitive.NewObjectID(),
		Identifier: identifier,
		Type:       logType,
		Message:    message,
		CreatedAt:  time.Now(),
	})
}
