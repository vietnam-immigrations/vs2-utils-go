package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colShopPrice = "price"

func CollectionPrice(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/db-shop")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colShopPrice)
}

type PriceKey string

const (
	PriceKeyEVisaNormal1MonthSingle PriceKey = "E_VISA_NORMAL_1_MONTH_SINGLE"
	PriceKeyEVisaNormal1MonthMulti  PriceKey = "E_VISA_NORMAL_1_MONTH_MULTI"
	PriceKeyEVisaNormal3MonthSingle PriceKey = "E_VISA_NORMAL_3_MONTH_SINGLE"
	PriceKeyEVisaNormal3MonthMulti  PriceKey = "E_VISA_NORMAL_3_MONTH_MULTI"
	PriceKeyEVisa2Day1MonthSingle   PriceKey = "E_VISA_2_DAY_1_MONTH_SINGLE"
	PriceKeyEVisa2Day1MonthMulti    PriceKey = "E_VISA_2_DAY_1_MONTH_MULTI"
	PriceKeyEVisa2Day3MonthSingle   PriceKey = "E_VISA_2_DAY_3_MONTH_SINGLE"
	PriceKeyEVisa2Day3MonthMulti    PriceKey = "E_VISA_2_DAY_3_MONTH_MULTI"
	PriceKeyEVisa1Day1MonthSingle   PriceKey = "E_VISA_1_DAY_1_MONTH_SINGLE"
	PriceKeyEVisa1Day1MonthMulti    PriceKey = "E_VISA_1_DAY_1_MONTH_MULTI"
	PriceKeyEVisa1Day3MonthSingle   PriceKey = "E_VISA_1_DAY_3_MONTH_SINGLE"
	PriceKeyEVisa1Day3MonthMulti    PriceKey = "E_VISA_1_DAY_3_MONTH_MULTI"
	PriceKeyEVisaUrgent1MonthSingle PriceKey = "E_VISA_URGENT_1_MONTH_SINGLE"
	PriceKeyEVisaUrgent1MonthMulti  PriceKey = "E_VISA_URGENT_1_MONTH_MULTI"
	PriceKeyEVisaUrgent3MonthSingle PriceKey = "E_VISA_URGENT_3_MONTH_SINGLE"
	PriceKeyEVisaUrgent3MonthMulti  PriceKey = "E_VISA_URGENT_3_MONTH_MULTI"
)

const (
	PriceKeyPriority2Day1MonthSingle   PriceKey = "PRIORITY_2_DAY_1_MONTH_SINGLE"
	PriceKeyPriority2Day1MonthMulti    PriceKey = "PRIORITY_2_DAY_1_MONTH_MULTI"
	PriceKeyPriority2Day3MonthSingle   PriceKey = "PRIORITY_2_DAY_3_MONTH_SINGLE"
	PriceKeyPriority2Day3MonthMulti    PriceKey = "PRIORITY_2_DAY_3_MONTH_MULTI"
	PriceKeyPriority1Day1MonthSingle   PriceKey = "PRIORITY_1_DAY_1_MONTH_SINGLE"
	PriceKeyPriority1Day1MonthMulti    PriceKey = "PRIORITY_1_DAY_1_MONTH_MULTI"
	PriceKeyPriority1Day3MonthSingle   PriceKey = "PRIORITY_1_DAY_3_MONTH_SINGLE"
	PriceKeyPriority1Day3MonthMulti    PriceKey = "PRIORITY_1_DAY_3_MONTH_MULTI"
	PriceKeyPriorityUrgent1MonthSingle PriceKey = "PRIORITY_URGENT_1_MONTH_SINGLE"
	PriceKeyPriorityUrgent1MonthMulti  PriceKey = "PRIORITY_URGENT_1_MONTH_MULTI"
	PriceKeyPriorityUrgent3MonthSingle PriceKey = "PRIORITY_URGENT_3_MONTH_SINGLE"
	PriceKeyPriorityUrgent3MonthMulti  PriceKey = "PRIORITY_URGENT_3_MONTH_MULTI"
)

const (
	PriceKeyVisaOnArrival1Day1MonthSingle PriceKey = "VISA_ON_ARRIVAL_1_DAY_1_MONTH_SINGLE"
)

const (
	PriceKeyFastTrackNormal PriceKey = "FAST_TRACK_NORMAL"
	PriceKeyFastTrackVIP    PriceKey = "FAST_TRACK_VIP"
	PriceKeyPickupCar       PriceKey = "PICKUP_CAR"
)

type Price struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Key   PriceKey           `bson:"key" json:"key"`
	Value string             `bson:"value" json:"value"`
}
