package push

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestSendToIOSApp(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	err := os.Setenv("APP", "vs2")
	assert.NoError(t, err)

	err = sendToIOSApp(ctx, "test title", "test description")
	assert.NoError(t, err)
}
