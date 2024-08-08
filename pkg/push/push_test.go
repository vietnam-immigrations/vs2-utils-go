package push_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/push"
)

func TestSend(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	err := os.Setenv("APP", "vs2")
	assert.NoError(t, err)
	push.SendNotificationForOrder(ctx, "testOrderID", "a title", "a message")
}
