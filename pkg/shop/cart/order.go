package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/random"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

func ToFinalOrder(ctx context.Context, uiOrder *db.UIOrder, prices []db.Price) (*db.Order, error) {
	log := logger.FromContext(ctx)
	log.Infof("Converting UI order to final order: %+v", *uiOrder)

	pm := newPriceManager(prices)

	finalOrder := &db.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: random.String(11, lo.NumbersCharset),
		CreatedAt:   time.Now(),
		Secret:      random.String(10, lo.AlphanumericCharset),
	}
	finalOrder.Applicants = uiOrder.Applicants
	finalOrder.PriorityApplicants = uiOrder.PriorityApplicants
	finalOrder.ApplicationType = uiOrder.ApplicationType
	finalOrder.Options = uiOrder.Options
	finalOrder.Billing = uiOrder.Billing
	finalOrder.BillingItems = make([]db.BillingItem, 0)

	noNormalApplicants := len(uiOrder.Applicants)
	noPriorityApplicants := len(uiOrder.PriorityApplicants)
	noApplicants := noNormalApplicants + noPriorityApplicants

	var billingVisaItem db.BillingItem
	if noPriorityApplicants > 0 {
		price, err := pm.GetPriorityPrice(uiOrder.Options.VisaType, uiOrder.Options.ProcessingTime)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get priority price")
		}
		priceInt := int(decimal.RequireFromString(price).IntPart())
		billingVisaItem = db.BillingItem{
			Description: fmt.Sprintf("[Priority] E-Visa %s - %s", uiOrder.Options.VisaType, uiOrder.Options.ProcessingTime),
			UnitPrice:   priceInt,
			Quantity:    noPriorityApplicants,
			Total:       priceInt * noPriorityApplicants,
		}
	}
	if noNormalApplicants > 0 {
		if uiOrder.ApplicationType == db.ApplicationTypeVisaOnArrival {
			price, err := pm.GetVOAPrice()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get VOA price")
			}
			priceInt := int(decimal.RequireFromString(price).IntPart())
			billingVisaItem = db.BillingItem{
				Description: fmt.Sprintf("Visa On Arrival %s", uiOrder.Options.VisaType),
				UnitPrice:   priceInt,
				Quantity:    noNormalApplicants,
				Total:       priceInt * noNormalApplicants,
			}
		} else {
			price, err := pm.GetEVisaPrice(uiOrder.Options.VisaType, uiOrder.Options.ProcessingTime)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get e-visa price")
			}
			priceInt := int(decimal.RequireFromString(price).IntPart())
			billingVisaItem = db.BillingItem{
				Description: fmt.Sprintf("E-Visa %s - %s", uiOrder.Options.VisaType, uiOrder.Options.ProcessingTime),
				UnitPrice:   priceInt,
				Quantity:    noNormalApplicants,
				Total:       priceInt * noNormalApplicants,
			}
		}
	}

	finalOrder.BillingItems = append(finalOrder.BillingItems, billingVisaItem)

	fastTrackPrice, err, ok := pm.GetFastTrackPrice(uiOrder.Options.FastTrack)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fast track price")
	}
	if ok {
		priceInt := int(decimal.RequireFromString(fastTrackPrice).IntPart())
		log.Infof("Adding fast track price: %+v", uiOrder.Options.FastTrack)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: uiOrder.Options.FastTrack,
			UnitPrice:   priceInt,
			Quantity:    noApplicants,
			Total:       priceInt * noApplicants,
		})
	}

	if uiOrder.Options.Car == db.CarYes {
		carPrice, err := pm.GetCarPrice()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get car price")
		}
		priceInt := int(decimal.RequireFromString(carPrice).IntPart())
		log.Infof("Adding car price: %+v", uiOrder.Options.Car)
		finalOrder.BillingItems = append(finalOrder.BillingItems, db.BillingItem{
			Description: "Car pickup",
			UnitPrice:   priceInt,
			Quantity:    1,
			Total:       priceInt,
		})
	}

	total := lo.SumBy(finalOrder.BillingItems, func(item db.BillingItem) int {
		return item.Total
	})
	finalOrder.Summary = db.OrderSummary{
		Total: total,
	}
	return finalOrder, nil
}
