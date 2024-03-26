package mail

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestSendUseBrevo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	_ = os.Setenv("APP", "vs2")
	_ = os.Setenv("LUG_LOCAL", "true")
	err := SendUseBrevo(ctx, ses.SendProps{
		From:    "info@vietnam-immigrations.org",
		To:      []string{"lenamtruong@gmail.com"},
		HTML:    "<h1>Test</h1>",
		Subject: "Test",
		CC:      []string{"namtruong.le@gmail.com"},
		BCC:     []string{"namtruongle2503@gmail.com"},
		ReplyTo: "info@vietnam-immigrations.org",
	})
	assert.NoError(t, err)
}
