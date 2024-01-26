package mail_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/mail"
)

func TestSendPriorityAdmin(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	_ = os.Setenv("APP", "vs2")
	_ = os.Setenv("LUG_LOCAL", "true")

	colOrders, err := db.CollectionOrders(ctx)
	assert.NoError(t, err)
	findOrder := colOrders.FindOne(ctx, bson.M{"number": "12621283177"})
	order := new(db.Order)
	err = findOrder.Decode(order)
	assert.NoError(t, err)
	err = mail.SendPriorityAdmin(ctx, order)
	assert.NoError(t, err)
}
