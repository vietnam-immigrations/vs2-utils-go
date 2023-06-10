package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

const (
	evisaPrice = 55
)

func ToFinalOrder(ctx context.Context, uiOrder *db.UIOrder) *db.Order {
	log := logger.FromContext(ctx)
	log.Infof("Converting UI order to final order: %+v", *uiOrder)
	finalOrder := &db.Order{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		Secret:    random.String(10, lo.AlphanumericCharset),
	}
	finalOrder.Applicants = uiOrder.Applicants
	finalOrder.Options = uiOrder.Options
	finalOrder.Billing = uiOrder.Billing
	numberOfApplicants := len(uiOrder.Applicants)
	finalOrder.BillingItems = []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   evisaPrice,
			Quantity:    numberOfApplicants,
			Total:       evisaPrice * numberOfApplicants,
		},
	}
	if processingTime, ok := db.ProcessingTimePrice[uiOrder.Options.ProcessingTime]; ok {
		log.Infof("Adding processing time price: %+v", uiOrder.Options.ProcessingTime)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: fmt.Sprintf("Processing time %s", uiOrder.Options.ProcessingTime),
			UnitPrice:   processingTime,
			Quantity:    numberOfApplicants,
			Total:       processingTime * numberOfApplicants,
		})
	}
	if fastTrack, ok := db.FastTrackPrice[uiOrder.Options.FastTrack]; ok {
		log.Infof("Adding fast track price: %+v", uiOrder.Options.FastTrack)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: uiOrder.Options.FastTrack,
			UnitPrice:   fastTrack,
			Quantity:    numberOfApplicants,
			Total:       fastTrack * numberOfApplicants,
		})
	}
	if car, ok := db.CarPrice[uiOrder.Options.Car]; ok {
		log.Infof("Adding car price: %+v", uiOrder.Options.Car)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: "Car pickup",
			UnitPrice:   car,
			Quantity:    1,
			Total:       car,
		})
	}

	total := lo.SumBy(finalOrder.BillingItems, func(item db.BillingItem) int {
		return item.Total
	})
	finalOrder.Summary = db.OrderSummary{
		Total: total,
	}
	return finalOrder
}
