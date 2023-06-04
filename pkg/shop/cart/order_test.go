package cart

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

func TestToFinalOrder(t *testing.T) {
	applicants := []db.Applicant{{
		PortraitFile: db.UploadFile{
			Name:   "portrait",
			Base64: "portrait-content",
		},
		PassportFile: db.UploadFile{
			Name:   "passport",
			Base64: "passport-content",
		},
		FirstName:      "Firstname",
		LastName:       "Lastname",
		DateOfBirth:    "01/01/1990",
		Sex:            "female",
		Nationality:    "Germany",
		PassportNumber: "DE999999",
		PassportExpiry: "01/01/2030",
	}, {
		PortraitFile: db.UploadFile{
			Name:   "portrait2",
			Base64: "portrait2-content",
		},
		PassportFile: db.UploadFile{
			Name:   "passport2",
			Base64: "passport2-content",
		},
		FirstName:      "Firstname2",
		LastName:       "Lastname2",
		DateOfBirth:    "01/01/2000",
		Sex:            "male",
		Nationality:    "USA",
		PassportNumber: "US1234567",
		PassportExpiry: "01/01/2040",
	}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTimeNormal,
		FastTrack:      db.FastTrackNo,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{
		FirstName: "Firstname",
		LastName:  "Lastname",
		Phone:     "+4912345678",
		Email:     "mail@mail.com",
		Email2:    "mail2@mail.com",
	}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.NotNil(t, finalOrder)
	assert.NotEmpty(t, finalOrder.ID)
	assert.Equal(t, applicants, finalOrder.Applicants)
	assert.Equal(t, options, finalOrder.Options)
	assert.Equal(t, billing, finalOrder.Billing)
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 110}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTime2Days(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTime2Days,
		FastTrack:      db.FastTrackNo,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time 2 working days",
			UnitPrice:   45,
			Quantity:    2,
			Total:       90,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 200}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTime1Day(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTime1Days,
		FastTrack:      db.FastTrackNo,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time 1 working day",
			UnitPrice:   60,
			Quantity:    2,
			Total:       120,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 230}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTimeSameDay(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTimeSameDay,
		FastTrack:      db.FastTrackNo,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time Same day",
			UnitPrice:   70,
			Quantity:    2,
			Total:       140,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 250}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTimeUrgent(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTimeUrgent,
		FastTrack:      db.FastTrackNo,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time Urgent",
			UnitPrice:   90,
			Quantity:    2,
			Total:       180,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 290}, finalOrder.Summary)
}

func TestToFinalOrder_FastTrackNormal(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTimeUrgent,
		FastTrack:      db.FastTrackNormal,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time Urgent",
			UnitPrice:   90,
			Quantity:    2,
			Total:       180,
		},
		{
			Description: "Normal fast-track",
			UnitPrice:   65,
			Quantity:    2,
			Total:       130,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 420}, finalOrder.Summary)
}

func TestToFinalOrder_FastTrackVIP(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTimeUrgent,
		FastTrack:      db.FastTrackVIP,
		Car:            db.CarNo,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time Urgent",
			UnitPrice:   90,
			Quantity:    2,
			Total:       180,
		},
		{
			Description: "VIP fast-track",
			UnitPrice:   95,
			Quantity:    2,
			Total:       190,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 480}, finalOrder.Summary)
}

func TestToFinalOrder_Car(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTimeUrgent,
		FastTrack:      db.FastTrackVIP,
		Car:            db.CarYes,
		Flight:         "VN123",
		Hotel:          "Hotel",
	}
	billing := db.CartBilling{}
	finalOrder := ToFinalOrder(context.TODO(), &db.UIOrder{
		Applicants: applicants,
		Options:    options,
		Billing:    billing,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
		{
			Description: "Processing time Urgent",
			UnitPrice:   90,
			Quantity:    2,
			Total:       180,
		},
		{
			Description: "VIP fast-track",
			UnitPrice:   95,
			Quantity:    2,
			Total:       190,
		},
		{
			Description: "Car pickup",
			UnitPrice:   35,
			Quantity:    1,
			Total:       35,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 515}, finalOrder.Summary)
}
