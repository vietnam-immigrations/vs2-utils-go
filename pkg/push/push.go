package push

import (
	"context"
	"fmt"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/logs"
)

func SendNotificationForOrder(ctx context.Context, orderID, title, message string) {
	log := logger.FromContext(ctx)

	// send ios push notification
	err := sendToIOSApp(ctx, title, message)
	if err != nil {
		log.Errorf("Failed to send to ios app: %s", err)
		// will not return here, continue to send pushover
	}

	logs.Add(ctx, db.LogTypeOrder, orderID, fmt.Sprintf("[%s] %s", title, message))
}
