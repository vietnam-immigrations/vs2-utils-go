package notification_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/notification"
)

func TestCreate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "prod")
	err := notification.Create(ctx, notification.Notification{
		ID:         uuid.New().String(),
		CSSClasses: "bg-positive text-white",
		Message:    "Test message",
	})
	assert.NoError(t, err)
}
