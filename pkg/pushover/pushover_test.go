package pushover_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/vietnam-immigrations/vs2-utils-go/pkg/pushover"
)

func TestSend(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	err := pushover.Send(ctx, "a title", "a message")
	assert.NoError(t, err)
}
