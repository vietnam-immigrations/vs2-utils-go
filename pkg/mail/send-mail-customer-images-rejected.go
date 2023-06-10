package mail

import (
	"context"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mail"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

func SendCustomerImagesRejected(ctx context.Context, order *db.Order) error {
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

	mjmlUsername, err := ssm.GetParameter(ctx, "/mjml/username", false)
	if err != nil {
		return err
	}
	mjmlPassword, err := ssm.GetParameter(ctx, "/mjml/password", true)
	if err != nil {
		return err
	}
	mailHTML, err := mail.Render(ctx, templateEmailImageRejected, templateEmailImageRejectedProps{
		OrderNumber: order.Number,
		UploadURL:   fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
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
		Subject: fmt.Sprintf("[IMPORTANT - PLEASE PROVIDE NEW IMAGES] Vietnam Visa Online Order #%s", order.Number),
		HTML:    *mailHTML,
	})
	return err
}
