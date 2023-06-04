package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colShopOrder = "order"

func CollectionOrder(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/mongo/db-shop", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colShopOrder)
}

type UploadFile struct {
	Name   string `bson:"name" json:"name"`
	Base64 string `bson:"base64" json:"base64"`
}

type Applicant struct {
	PortraitFile   UploadFile `bson:"portraitFile" json:"portraitFile"`
	PassportFile   UploadFile `bson:"passportFile" json:"passportFile"`
	FirstName      string     `bson:"firstName" json:"firstName"`
	LastName       string     `bson:"lastName" json:"lastName"`
	DateOfBirth    string     `bson:"dateOfBirth" json:"dateOfBirth"`
	Sex            string     `bson:"sex" json:"sex"`
	Nationality    string     `bson:"nationality" json:"nationality"`
	PassportNumber string     `bson:"passportNumber" json:"passportNumber"`
	PassportExpiry string     `bson:"passportExpiry" json:"passportExpiry"`
}

type CartOptions struct {
	ArrivalDate    string `bson:"arrivalDate" json:"arrivalDate"`
	Entry          string `bson:"entry" json:"entry"`
	ProcessingTime string `bson:"processingTime" json:"processingTime"`
	FastTrack      string `bson:"fastTrack" json:"fastTrack"`
	Car            string `bson:"car" json:"car"`
	Flight         string `bson:"flight" json:"flight"`
	Hotel          string `bson:"hotel" json:"hotel"`
}

type CartBilling struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Phone     string `bson:"phone" json:"phone"`
	Email     string `bson:"email" json:"email"`
	Email2    string `bson:"email2" json:"email2"`
}

type UIOrder struct {
	Applicants []Applicant `bson:"applicants" json:"applicants"`
	Options    CartOptions `bson:"options" json:"options"`
	Billing    CartBilling `bson:"billing" json:"billing"`
}

type BillingItem struct {
	Description string `bson:"description" json:"description"`
	UnitPrice   int    `bson:"unitPrice" json:"unitPrice"`
	Quantity    int    `bson:"quantity" json:"quantity"`
	Total       int    `bson:"total" json:"total"`
}

type OrderSummary struct {
	Total int `bson:"total" json:"total"`
}

type Order struct {
	UIOrder
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	BillingItems []BillingItem      `bson:"billingItems" json:"billingItems"`
	Summary      OrderSummary       `bson:"summary" json:"summary"`
}

const (
	ProcessingTimeNormal  = "10-12 working days"
	ProcessingTime2Days   = "2 working days"
	ProcessingTime1Days   = "1 working day"
	ProcessingTimeSameDay = "Same day"
	ProcessingTimeUrgent  = "Urgent"
)

var ProcessingTimePrice = map[string]int{
	ProcessingTime2Days:   45,
	ProcessingTime1Days:   60,
	ProcessingTimeSameDay: 70,
	ProcessingTimeUrgent:  90,
}

const (
	FastTrackNo     = "No"
	FastTrackNormal = "Normal fast-track"
	FastTrackVIP    = "VIP fast-track"
)

var FastTrackPrice = map[string]int{
	FastTrackNormal: 65,
	FastTrackVIP:    95,
}

const (
	CarNo  = "No"
	CarYes = "Yes"
)

var CarPrice = map[string]int{
	CarYes: 35,
}
