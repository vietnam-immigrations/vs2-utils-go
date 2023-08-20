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
	noPriorityApplicants := len(uiOrder.PriorityApplicants)
	noApplicants := noNormalApplicants + noPriorityApplicants

	var billingVisaItem db.BillingItem
	if noPriorityApplicants > 0 {
		billingVisaItem = db.BillingItem{
			Description: fmt.Sprintf("[Priority] E-Visa %s", uiOrder.Options.VisaType),
			UnitPrice:   VisaPricePriority[uiOrder.Options.VisaType],
			Quantity:    noPriorityApplicants,
			Total:       VisaPricePriority[uiOrder.Options.VisaType] * noPriorityApplicants,
		}
	}
	if noNormalApplicants > 0 {
		billingVisaItem = db.BillingItem{
			Description: fmt.Sprintf("E-Visa %s", uiOrder.Options.VisaType),
			UnitPrice:   VisaPriceStandard[uiOrder.Options.VisaType],
			Quantity:    noNormalApplicants,
			Total:       VisaPriceStandard[uiOrder.Options.VisaType] * noNormalApplicants,
		}
	}

	if processingTime, ok := ProcessingTimePrice[uiOrder.Options.VisaType][uiOrder.Options.ProcessingTime]; ok {
		log.Infof("Adding processing time price: %+v", uiOrder.Options.ProcessingTime)
		billingVisaItem.Description = fmt.Sprintf("%s - %s", billingVisaItem.Description, uiOrder.Options.ProcessingTime)
		billingVisaItem.UnitPrice = processingTime + billingVisaItem.UnitPrice
		billingVisaItem.Total = noApplicants * billingVisaItem.UnitPrice
	}
	finalOrder.BillingItems = append(finalOrder.BillingItems, billingVisaItem)

	if fastTrack, ok := FastTrackPrice[uiOrder.Options.FastTrack]; ok {
		log.Infof("Adding fast track price: %+v", uiOrder.Options.FastTrack)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: uiOrder.Options.FastTrack,
			UnitPrice:   fastTrack,
			Quantity:    noApplicants,
			Total:       fastTrack * noApplicants,
		})
	}
	if car, ok := CarPrice[uiOrder.Options.Car]; ok {
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
