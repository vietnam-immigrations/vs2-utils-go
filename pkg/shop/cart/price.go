package cart

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

type priceManager struct {
	prices []db.Price
}

func newPriceManager(prices []db.Price) *priceManager {
	return &priceManager{
		prices: prices,
	}
}

func (pm *priceManager) GetPrice(key db.PriceKey) (string, error) {
	price, ok := lo.Find(pm.prices, func(price db.Price) bool {
		return price.Key == key
	})
	if !ok {
		return "", fmt.Errorf("price not found for key %s", key)
	}
	return price.Value, nil
}

func (pm *priceManager) GetVOAPrice() (string, error) {
	return pm.GetPrice(db.PriceKeyVisaOnArrival1Day1MonthSingle)
}

func (pm *priceManager) GetFastTrackPrice(ft string) (string, error, bool) {
	switch ft {
	case db.FastTrackNormal:
		val, err := pm.GetPrice(db.PriceKeyFastTrackNormal)
		return val, err, true
	case db.FastTrackVIP:
		val, err := pm.GetPrice(db.PriceKeyFastTrackVIP)
		return val, err, true
	default:
		return "", nil, false
	}
}

func (pm *priceManager) GetCarPrice() (string, error) {
	return pm.GetPrice(db.PriceKeyPickupCar)
}

func (pm *priceManager) GetEVisaNormalPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisaNormal1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisaNormal1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisaNormal3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisaNormal3MonthMulti)
	default:
		return "", fmt.Errorf("no evisa normal price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetEVisa2DayPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisa2Day1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisa2Day1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisa2Day3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisa2Day3MonthMulti)
	default:
		return "", fmt.Errorf("no evisa 2 day price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetEVisa1DayPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisa1Day1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisa1Day1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisa1Day3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisa1Day3MonthMulti)
	default:
		return "", fmt.Errorf("no evisa 1 day price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetEVisaUrgentPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisaUrgent1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisaUrgent1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyEVisaUrgent3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyEVisaUrgent3MonthMulti)
	default:
		return "", fmt.Errorf("no evisa urgent price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetEVisaPrice(visaType string, processingTime string) (string, error) {
	switch processingTime {
	case db.ProcessingTimeNormal:
		return pm.GetEVisaNormalPrice(visaType)
	case db.ProcessingTime2Days:
		return pm.GetEVisa2DayPrice(visaType)
	case db.ProcessingTime1Days:
		return pm.GetEVisa1DayPrice(visaType)
	case db.ProcessingTimeUrgent:
		return pm.GetEVisaUrgentPrice(visaType)
	default:
		return "", fmt.Errorf("no evisa price configured for processing time %s, visa type %s", processingTime, visaType)
	}
}

func (pm *priceManager) GetPriority2DayPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyPriority2Day1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyPriority2Day1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyPriority2Day3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyPriority2Day3MonthMulti)
	default:
		return "", fmt.Errorf("no priority 2 day price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetPriority1DayPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyPriority1Day1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyPriority1Day1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyPriority1Day3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyPriority1Day3MonthMulti)
	default:
		return "", fmt.Errorf("no priority 1 day price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetPriorityUrgentPrice(visaType string) (string, error) {
	switch visaType {
	case db.VisaType1MonthSingle:
		return pm.GetPrice(db.PriceKeyPriorityUrgent1MonthSingle)
	case db.VisaType1MonthMultiple:
		return pm.GetPrice(db.PriceKeyPriorityUrgent1MonthMulti)
	case db.VisaType3MonthSingle:
		return pm.GetPrice(db.PriceKeyPriorityUrgent3MonthSingle)
	case db.VisaType3MonthMultiple:
		return pm.GetPrice(db.PriceKeyPriorityUrgent3MonthMulti)
	default:
		return "", fmt.Errorf("no priority urgent price configured for visa type %s", visaType)
	}
}

func (pm *priceManager) GetPriorityPrice(visaType string, processingTime string) (string, error) {
	switch processingTime {
	case db.ProcessingTime2Days:
		return pm.GetPriority2DayPrice(visaType)
	case db.ProcessingTime1Days:
		return pm.GetPriority1DayPrice(visaType)
	case db.ProcessingTimeUrgent:
		return pm.GetPriorityUrgentPrice(visaType)
	default:
		return "", fmt.Errorf("no priority price configured for processing time %s, visa type %s", processingTime, visaType)
	}
}
