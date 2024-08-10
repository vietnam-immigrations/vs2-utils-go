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
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/push"
)

func SendPriorityCustomer(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send priority email to customer for order [%s]", order.Number)

	cfg, err := db.GetConfig(ctx)
	if err != nil {
		return err
	}

	mjmlUsername, err := secretsmanager.GetParameter(ctx, "/mjml/username")
	if err != nil {
		return err
	}
	mjmlPassword, err := secretsmanager.GetParameter(ctx, "/mjml/password")
	if err != nil {
		return err
	}
	mailHTML, err := mail.Render(ctx, templateEmailPriorityCustomer, templateEmailPriorityCustomerProps{
		FullName:    strings.ToUpper(fmt.Sprintf("%s %s", order.Billing.FirstName, order.Billing.LastName)),
		ArrivalDate: order.Trip.ArrivalDate,
		Entry:       order.Trip.Checkpoint,
		Applicants: lo.Map(order.Applicants, func(app db.Applicant, i int) templateEmailPriorityCustomerPropsApplicant {
			return templateEmailPriorityCustomerPropsApplicant{
				Title:            fmt.Sprintf("Applicant %d", i+1),
				RegistrationCode: app.RegistrationCode,
			}
		}),
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}
	err = SendUseBrevo(ctx, ses.SendProps{
		From:        mailAddressInfo,
		To:          lo.Compact([]string{order.Billing.Email, order.Billing.Email2}),
		ReplyTo:     mailAddressInfo,
		BCC:         nil,
		CC:          []string{cfg.EmailCustomerCC},
		Subject:     fmt.Sprintf("Vietnam Visa Priority Handling #%s", order.Number),
		HTML:        *mailHTML,
		Attachments: nil,
	})

	if err != nil {
		push.SendNotificationForOrder(
			ctx,
			order.ID.Hex(),
			"Failed to send priority email to customer",
			fmt.Sprintf("Failed to send priority email to customer for order [%s]: %s", order.Number, err),
		)
		return err
	}

	push.SendNotificationForOrder(
		ctx,
		order.ID.Hex(),
		"Priority email sent to customer",
		fmt.Sprintf("Priority email sent to customer for order [%s]", order.Number),
	)

	return nil
}
