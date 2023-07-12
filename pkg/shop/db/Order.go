package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colShopOrder = "order"

func CollectionOrder(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/db-shop")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colShopOrder)
}

type Applicant struct {
	PortraitFile   string `bson:"portraitFile" json:"portraitFile"`
	PassportFile   string `bson:"passportFile" json:"passportFile"`
	FirstName      string `bson:"firstName" json:"firstName"`
	LastName       string `bson:"lastName" json:"lastName"`
	DateOfBirth    string `bson:"dateOfBirth" json:"dateOfBirth"`
	Sex            string `bson:"sex" json:"sex"`
	Nationality    string `bson:"nationality" json:"nationality"`
	PassportNumber string `bson:"passportNumber" json:"passportNumber"`
	PassportExpiry string `bson:"passportExpiry" json:"passportExpiry"`
}

type PriorityApplicant struct {
	Code  string `bson:"code" json:"code"`
	Email string `bson:"email" json:"email"`
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
	Applicants         []Applicant         `bson:"applicants" json:"applicants"`
	PriorityApplicants []PriorityApplicant `bson:"priorityApplicants" json:"priorityApplicants"`
	Options            CartOptions         `bson:"options" json:"options"`
	Billing            CartBilling         `bson:"billing" json:"billing"`
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
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	OrderNumber      string             `bson:"orderNumber" json:"orderNumber"`
	BillingItems     []BillingItem      `bson:"billingItems" json:"billingItems"`
	Summary          OrderSummary       `bson:"summary" json:"summary"`
	Secret           string             `bson:"secret" json:"secret"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	PaymentCreatedAt *time.Time         `bson:"paymentCreatedAt,omitempty" json:"paymentCreatedAt,omitempty"`
	PaidAt           *time.Time         `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
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
