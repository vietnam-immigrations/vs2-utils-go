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

func TestSendAdmin(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "prod")
	_ = os.Setenv("APP", "vs2")
	_ = os.Setenv("LUG_LOCAL", "true")

	colOrders, err := db.CollectionOrders(ctx)
	assert.NoError(t, err)
	findOrder := colOrders.FindOne(ctx, bson.M{"number": "50660878185"})
	order := new(db.Order)
	err = findOrder.Decode(order)
	assert.NoError(t, err)
	err = mail.SendAdmin(ctx, order)
	assert.NoError(t, err)
}
