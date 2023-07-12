package mail

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mail"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

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

	mjmlUsername, err := secretsmanager.GetParameter(ctx, "/mjml/username")
	if err != nil {
		return err
	}
	mjmlPassword, err := secretsmanager.GetParameter(ctx, "/mjml/password")
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
		From: mailAddressInfo,
		To: lo.Compact([]string{
			order.Billing.Email, order.Billing.Email2,
		}),
		CC:      []string{cfg.EmailCustomerCC},
		ReplyTo: mailAddressInfo,
		Subject: fmt.Sprintf("Vietnam Visa Online Order #%s", order.Number),
		HTML:    *mailHTML,
	})
	return err
}
