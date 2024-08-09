package mail_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/mail"
	shopdb "github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

func TestSendPaymentPending(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	_ = os.Setenv("APP", "vs2")
	_ = os.Setenv("LUG_LOCAL", "true")

	colOrders, err := shopdb.CollectionOrder(ctx)
	assert.NoError(t, err)
	findOrder := colOrders.FindOne(ctx, bson.M{"orderNumber": "34183804613"})
	order := new(shopdb.Order)
	err = findOrder.Decode(order)
	assert.NoError(t, err)
	err = mail.SendPaymentPending(ctx, order)
	assert.NoError(t, err)
}
