package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/pkg/mongodb"
)

const (
	colOrdersName  = "orders"
	keyOrderID     = "order_id"
	keyOrderWooID  = "order_woo_id"
	keyOrderNumber = "order_number"
)

func CollectionOrders(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colOrdersName)
}

// AddOrderToContext adds order data to context
func AddOrderToContext(ctx context.Context, order Order) context.Context {
	logger.AddFields(keyOrderID, keyOrderWooID, keyOrderNumber)
	result := context.WithValue(ctx, keyOrderID, order.ID)
	result = context.WithValue(result, keyOrderWooID, order.OrderID)
	result = context.WithValue(result, keyOrderNumber, order.Number)
	return result
}

type Billing struct {
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Email     string `bson:"email" json:"email"`
	Email2    string `bson:"email2" json:"email2"`
	Phone     string `bson:"phone" json:"phone"`
}

type Trip struct {
	// ArrivalDate is the raw value from woocommerce order
	ArrivalDate      string `bson:"arrivalDate" json:"arrivalDate"`
	Checkpoint       string `bson:"checkpoint" json:"checkpoint"`
	ProcessingTime   string `bson:"processingTime" json:"processingTime"`
	FastTrack        string `bson:"fastTrack" json:"fastTrack"`
	CarPickup        bool   `bson:"carPickup" json:"carPickup"`
	Flight           string `bson:"flight" json:"flight"`
	CarPickupAddress string `bson:"carPickupAddress" json:"carPickupAddress"`

	// Arrival is the parsed value from ArrivalDate
	Arrival time.Time `bson:"arrival" json:"arrival"`
}

type Applicant struct {
	FirstName        string `bson:"firstName" json:"firstName"`
	LastName         string `bson:"lastName" json:"lastName"`
	DateOfBirth      string `bson:"dateOfBirth" json:"dateOfBirth"`
	Sex              string `bson:"sex" json:"sex"`
	Nationality      string `bson:"nationality" json:"nationality"`
	PassportNumber   string `bson:"passportNumber" json:"passportNumber"`
	PassportExpiry   string `bson:"passportExpiry" json:"passportExpiry"`
	RegistrationCode string `bson:"registrationCode" json:"registrationCode"`
	Email            string `bson:"email" json:"email"`

	VisaS3Key    string `bson:"visaS3Key" json:"visaS3Key"`
	VisaSent     bool   `bson:"visaSent" json:"visaSent"`
	CancelReason string `bson:"cancelReason" json:"cancelReason"`
}

type OrderType string

const (
	OrderTypeVisa     OrderType = ""
	OrderTypePriority OrderType = "Priority"
)

type Order struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	OrderID            int                `bson:"orderId" json:"orderId"`
	Total              string             `bson:"total" json:"total"`
	OrderKey           string             `bson:"orderKey" json:"orderKey"`
	Billing            Billing            `bson:"billing" json:"billing"`
	PaymentMethodTitle string             `bson:"paymentMethodTitle" json:"paymentMethodTitle"`
	Number             string             `bson:"number" json:"number"`
	Trip               Trip               `bson:"trip" json:"trip"`
	Applicants         []Applicant        `bson:"applicants" json:"applicants"`
	Type               OrderType          `bson:"type" json:"type"`

	AllVisaSent  bool      `bson:"allVisaSent" json:"allVisaSent"`
	InvoiceDocID string    `bson:"invoiceDocId" json:"invoiceDocId"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
}
