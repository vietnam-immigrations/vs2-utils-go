package cart

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

func TestToFinalOrder(t *testing.T) {
	applicants := []db.Applicant{{
		PortraitFile:   "portrait-file",
		PassportFile:   "passport-file",
		FirstName:      "Firstname",
		LastName:       "Lastname",
		DateOfBirth:    "01/01/1990",
		Sex:            "female",
		Nationality:    "Germany",
		PassportNumber: "DE999999",
		PassportExpiry: "01/01/2030",
	}, {
		PortraitFile:   "portrait-file-2",
		PassportFile:   "passport-file-2",
		FirstName:      "Firstname2",
		LastName:       "Lastname2",
		DateOfBirth:    "01/01/2000",
		Sex:            "male",
		Nationality:    "USA",
		PassportNumber: "US1234567",
		PassportExpiry: "01/01/2040",
	}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
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
		Applicants:      applicants,
		Options:         options,
		Billing:         billing,
		ApplicationType: db.ApplicationTypeEVisa,
	})
	assert.NotNil(t, finalOrder)
	assert.Equal(t, finalOrder.UIOrder.ApplicationType, db.ApplicationTypeEVisa)
	assert.NotEmpty(t, finalOrder.ID)
	assert.NotEmpty(t, finalOrder.Secret)
	assert.NotEmpty(t, finalOrder.CreatedAt)
	assert.Nil(t, finalOrder.PaidAt)
	assert.Equal(t, applicants, finalOrder.Applicants)
	assert.Equal(t, options, finalOrder.Options)
	assert.Equal(t, billing, finalOrder.Billing)
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa 1 month single entry",
			UnitPrice:   55,
			Quantity:    2,
			Total:       110,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 110}, finalOrder.Summary)
	assert.Equal(t, len(finalOrder.OrderNumber), 11)
}

func TestToFinalOrder_3Month(t *testing.T) {
	applicants := []db.Applicant{{
		PortraitFile:   "portrait-file",
		PassportFile:   "passport-file",
		FirstName:      "Firstname",
		LastName:       "Lastname",
		DateOfBirth:    "01/01/1990",
		Sex:            "female",
		Nationality:    "Germany",
		PassportNumber: "DE999999",
		PassportExpiry: "01/01/2030",
	}, {
		PortraitFile:   "portrait-file-2",
		PassportFile:   "passport-file-2",
		FirstName:      "Firstname2",
		LastName:       "Lastname2",
		DateOfBirth:    "01/01/2000",
		Sex:            "male",
		Nationality:    "USA",
		PassportNumber: "US1234567",
		PassportExpiry: "01/01/2040",
	}}
	options := db.CartOptions{
		VisaType:       db.VisaType3MonthSingle,
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
	assert.NotEmpty(t, finalOrder.Secret)
	assert.NotEmpty(t, finalOrder.CreatedAt)
	assert.Nil(t, finalOrder.PaidAt)
	assert.Equal(t, applicants, finalOrder.Applicants)
	assert.Equal(t, options, finalOrder.Options)
	assert.Equal(t, billing, finalOrder.Billing)
	assert.Equal(t, []db.BillingItem{
		{
			Description: "E-Visa 3 months single entry",
			UnitPrice:   116,
			Quantity:    2,
			Total:       232,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 232}, finalOrder.Summary)
	assert.Equal(t, len(finalOrder.OrderNumber), 11)
}

func TestToFinalOrder_Priority(t *testing.T) {
	priorityApplicants := []db.PriorityApplicant{{
		Code:  "code1",
		Email: "email1",
	}, {
		Code:  "code2",
		Email: "email2",
	}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
		ArrivalDate:    "01/01/2023",
		Entry:          "Hanoi",
		ProcessingTime: db.ProcessingTime2Days,
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
		Applicants:         make([]db.Applicant, 0),
		PriorityApplicants: priorityApplicants,
		Options:            options,
		Billing:            billing,
	})
	assert.NotNil(t, finalOrder)
	assert.NotEmpty(t, finalOrder.ID)
	assert.NotEmpty(t, finalOrder.Secret)
	assert.NotEmpty(t, finalOrder.CreatedAt)
	assert.Nil(t, finalOrder.PaidAt)
	assert.Equal(t, priorityApplicants, finalOrder.PriorityApplicants)
	assert.Equal(t, options, finalOrder.Options)
	assert.Equal(t, billing, finalOrder.Billing)
	assert.Equal(t, []db.BillingItem{
		{
			Description: "[Priority] E-Visa 1 month single entry - 2 working days",
			UnitPrice:   75,
			Quantity:    2,
			Total:       150,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 150}, finalOrder.Summary)
	assert.Equal(t, len(finalOrder.OrderNumber), 11)
}

func TestToFinalOrder_ProcessingTime2Days(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
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
			Description: "E-Visa 1 month single entry - 2 working days",
			UnitPrice:   100,
			Quantity:    2,
			Total:       200,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 200}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTime1Day(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
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
			Description: "E-Visa 1 month single entry - 1 working day",
			UnitPrice:   115,
			Quantity:    2,
			Total:       230,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 230}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTime1Day_VOA(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
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
		Applicants:      applicants,
		Options:         options,
		Billing:         billing,
		ApplicationType: db.ApplicationTypeVisaOnArrival,
	})
	assert.Equal(t, []db.BillingItem{
		{
			Description: "Visa On Arrival 1 month single entry",
			UnitPrice:   125,
			Quantity:    2,
			Total:       250,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 250}, finalOrder.Summary)
}

func TestToFinalOrder_ProcessingTimeUrgent(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
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
			Description: "E-Visa 1 month single entry - Urgent",
			UnitPrice:   145,
			Quantity:    2,
			Total:       290,
		},
	}, finalOrder.BillingItems)
	assert.Equal(t, db.OrderSummary{Total: 290}, finalOrder.Summary)
}

func TestToFinalOrder_FastTrackNormal(t *testing.T) {
	applicants := []db.Applicant{{}, {}}
	options := db.CartOptions{
		VisaType:       db.VisaType1MonthSingle,
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
			Description: "E-Visa 1 month single entry - Urgent",
			UnitPrice:   145,
			Quantity:    2,
			Total:       290,
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
		VisaType:       db.VisaType1MonthSingle,
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
			Description: "E-Visa 1 month single entry - Urgent",
			UnitPrice:   145,
			Quantity:    2,
			Total:       290,
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
		VisaType:       db.VisaType1MonthSingle,
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
			Description: "E-Visa 1 month single entry - Urgent",
			UnitPrice:   145,
			Quantity:    2,
			Total:       290,
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
