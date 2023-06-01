package mail

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mail"
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

	// extra services
	extraServices := make([]string, 0)
	if order.Trip.FastTrack != "No" {
		extraServices = append(extraServices, fmt.Sprintf("%s (flight: %s)", order.Trip.FastTrack, order.Trip.Flight))
	}
	if order.Trip.CarPickup {
		extraServices = append(extraServices, fmt.Sprintf("car pick-up (hotel: %s)", order.Trip.CarPickupAddress))
	}
	extraServicesText := strings.Join(extraServices, ", ")

	mjmlUsername, err := ssm.GetParameter(ctx, "/mjml/username", false)
	if err != nil {
		return err
	}
	mjmlPassword, err := ssm.GetParameter(ctx, "/mjml/password", true)
	if err != nil {
		return err
	}
	mailHTML, err := mail.Render(ctx, templateEmailCustomer, templateEmailCustomerProps{
		OrderNumber:      order.Number,
		ArrivalDate:      order.Trip.ArrivalDate,
		Entry:            order.Trip.Checkpoint,
		Flight:           order.Trip.Flight,
		Hotel:            order.Trip.CarPickupAddress,
		ProcessingTime:   order.Trip.ProcessingTime,
		Telephone:        order.Billing.Phone,
		Email:            order.Billing.Email,
		Email2:           order.Billing.Email2,
		HasExtraServices: len(extraServices) > 0,
		ExtraServices:    extraServicesText,
		Applicants: lo.Map(order.Applicants, func(app db.Applicant, _ int) templateEmailCustomerPropsApplicant {
			return templateEmailCustomerPropsApplicant{
				Name:               fmt.Sprintf("%s, %s", app.LastName, app.FirstName),
				Nationality:        app.Nationality,
				Passport:           app.PassportNumber,
				Birthday:           app.DateOfBirth,
				PassportValidUntil: app.PassportExpiry,
				Gender:             app.Sex,
			}
		}),
		TrackingURL: fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}
	err = ses.Send(ctx, ses.SendProps{
		From: "info@vietnam-immigrations.org",
		To: lo.Compact([]string{
			"info@vietnam-immigrations.org",
			order.Billing.Email, order.Billing.Email2,
		}),
		ReplyTo: "info@vietnam-immigrations.org",
		Subject: fmt.Sprintf("Vietnam Visa Online Order #%s", order.Number),
		HTML:    *mailHTML,
	})
	return err
}
