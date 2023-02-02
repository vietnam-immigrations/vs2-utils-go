package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	mymailjet "github.com/nam-truong-le/lambda-utils-go/v3/pkg/mailjet"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

type SendCustomerOptionsApplicant struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Sex                string `json:"sex"`
	DateOfBirth        string `json:"dateOfBirth"`
	Nationality        string `json:"nationality"`
	PassportNumber     string `json:"passportNumber"`
	PassportValidUntil string `json:"passportValidUntil"`
}

type SendCustomerOptions struct {
	Applicants     []SendCustomerOptionsApplicant `json:"applicants"`
	OrderNumber    string                         `json:"orderNumber"`
	ArrivalDate    string                         `json:"arrivalDate"`
	Entry          string                         `json:"entry"`
	Flight         string                         `json:"flight"`
	Hotel          string                         `json:"hotel"`
	ProcessingTime string                         `json:"processingTime"`
	Telephone      string                         `json:"telephone"`
	Email          string                         `json:"email"`
	SecondaryEmail string                         `json:"secondaryEmail"`
	ExtraServices  string                         `json:"extraServices"`
	StatusUrl      string                         `json:"statusUrl"`
}

func SendCustomer(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send email to customer for order [%s]", order.Number)

	cfg, err := db.GetConfig(ctx)
	if err != nil {
		return err
	}

	to := mailjet.RecipientsV31{
		{
			Email: order.Billing.Email,
			Name:  fmt.Sprintf("%s %s", order.Billing.LastName, order.Billing.FirstName),
		},
	}
	if order.Billing.Email2 != "" {
		to = append(to, mailjet.RecipientV31{
			Email: order.Billing.Email2,
			Name:  fmt.Sprintf("%s %s", order.Billing.LastName, order.Billing.FirstName),
		})
	}

	applicants := make([]SendCustomerOptionsApplicant, 0)
	for _, app := range order.Applicants {
		applicants = append(applicants, SendCustomerOptionsApplicant{
			FirstName:          app.FirstName,
			LastName:           app.LastName,
			Sex:                app.Sex,
			DateOfBirth:        app.DateOfBirth,
			Nationality:        app.Nationality,
			PassportNumber:     app.PassportNumber,
			PassportValidUntil: app.PassportExpiry,
		})
	}

	// extra services
	ft := make([]string, 0)
	if order.Trip.FastTrack != "No" {
		ft = append(ft, fmt.Sprintf("%s (flight: %s)", order.Trip.FastTrack, order.Trip.Flight))
	}
	if order.Trip.CarPickup {
		ft = append(ft, fmt.Sprintf("car pick-up (hotel: %s)", order.Trip.CarPickupAddress))
	}
	ftText := strings.Join(ft, ", ")

	variables := SendCustomerOptions{
		Applicants:     applicants,
		OrderNumber:    order.Number,
		ArrivalDate:    order.Trip.ArrivalDate,
		Entry:          order.Trip.Checkpoint,
		Flight:         order.Trip.Flight,
		Hotel:          order.Trip.CarPickupAddress,
		ProcessingTime: order.Trip.ProcessingTime,
		Telephone:      order.Billing.Phone,
		Email:          order.Billing.Email,
		SecondaryEmail: order.Billing.Email2,
		ExtraServices:  ftText,
		StatusUrl:      fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
	}
	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		return nil
	}
	rawVariables := new(map[string]interface{})
	err = json.Unmarshal(jsonVariables, rawVariables)
	if err != nil {
		return nil
	}

	log.Infof("%+v", variables)

	body := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: "info@vietnam-immigrations.org",
			Name:  "Vietnam Visa Online",
		},
		To: &to,
		Cc: &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: cfg.EmailCustomerCC,
			},
		},
		TemplateID:       cfg.EmailCustomerTemplateID,
		TemplateLanguage: true,
		Subject:          fmt.Sprintf("Vietnam Visa Online Order #%s", order.Number),
		Variables:        *rawVariables,
	}

	return mymailjet.Send(ctx, body)
}
