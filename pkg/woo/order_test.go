package woo

import (
	"context"
	"testing"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestGetOrder(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	_, err := GetOrder(ctx, "110")
	if err != nil {
		t.Fail()
	}
}
