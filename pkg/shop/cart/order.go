package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/random"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

const (
	evisaPrice    = 55
	priorityPrice = 25
)

func ToFinalOrder(ctx context.Context, uiOrder *db.UIOrder) *db.Order {
	log := logger.FromContext(ctx)
	log.Infof("Converting UI order to final order: %+v", *uiOrder)
	finalOrder := &db.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: random.String(11, lo.NumbersCharset),
		CreatedAt:   time.Now(),
		Secret:      random.String(10, lo.AlphanumericCharset),
	}
	finalOrder.Applicants = uiOrder.Applicants
	finalOrder.PriorityApplicants = uiOrder.PriorityApplicants
	finalOrder.Options = uiOrder.Options
	finalOrder.Billing = uiOrder.Billing
	finalOrder.BillingItems = make([]db.BillingItem, 0)
	noNormalApplicants := len(uiOrder.Applicants)
	if noNormalApplicants > 0 {
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: "E-Visa",
			UnitPrice:   evisaPrice,
			Quantity:    noNormalApplicants,
			Total:       evisaPrice * noNormalApplicants,
		})
	}
	noPriorityApplicants := len(uiOrder.PriorityApplicants)
	if noPriorityApplicants > 0 {
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: "E-Visa Priority Handling",
			UnitPrice:   priorityPrice,
			Quantity:    noPriorityApplicants,
			Total:       priorityPrice * noPriorityApplicants,
		})
	}
	noApplicants := noNormalApplicants + noPriorityApplicants
	if processingTime, ok := db.ProcessingTimePrice[uiOrder.Options.ProcessingTime]; ok {
		log.Infof("Adding processing time price: %+v", uiOrder.Options.ProcessingTime)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: fmt.Sprintf("Processing time %s", uiOrder.Options.ProcessingTime),
			UnitPrice:   processingTime,
			Quantity:    noApplicants,
			Total:       processingTime * noApplicants,
		})
	}
	if fastTrack, ok := db.FastTrackPrice[uiOrder.Options.FastTrack]; ok {
		log.Infof("Adding fast track price: %+v", uiOrder.Options.FastTrack)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: uiOrder.Options.FastTrack,
			UnitPrice:   fastTrack,
			Quantity:    noApplicants,
			Total:       fastTrack * noApplicants,
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
