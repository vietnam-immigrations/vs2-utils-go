package mail

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/push"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mail"
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
		VisaType:         order.Trip.VisaType,
		VisitPurpose:     order.Trip.VisitPurpose,
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
				Gender:             app.Sex,
				Nationality:        app.Nationality,
				Passport:           app.PassportNumber,
				Birthday:           app.DateOfBirth,
				PassportValidUntil: app.PassportExpiry,
				HomeAddress:        app.AddressHome,
				HomeContact:        app.PhoneNumberHome,
				VietnamAddress:     app.AddressVietnam,
				PreviousVisitCount: app.PreviousVisitCount,
				LawViolation:       app.LawViolation,
			}
		}),
		TrackingURL:     fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
		IsVisaOnArrival: order.Variant == db.OrderVariantVisaOnArrival,
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}
	err = SendUseBrevo(ctx, ses.SendProps{
		From: mailAddressInfo,
		To: lo.Compact([]string{
			order.Billing.Email, order.Billing.Email2,
		}),
		CC:      []string{cfg.EmailCustomerCC},
		ReplyTo: mailAddressInfo,
		Subject: fmt.Sprintf("Vietnam Visa Online Order #%s", order.Number),
		HTML:    *mailHTML,
	})

	if err != nil {
		push.SendNotificationForOrder(
			ctx,
			order.ID.Hex(),
			"Failed to send email to customer",
			fmt.Sprintf("Failed to send email to customer for order [%s]: %s", order.Number, err),
		)
		return err
	}

	push.SendNotificationForOrder(
		ctx,
		order.ID.Hex(),
		"Email sent to customer",
		fmt.Sprintf("Email sent to customer for order [%s]", order.Number),
	)

	return nil
}
