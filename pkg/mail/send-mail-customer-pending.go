package mail

import (
	"context"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mail"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

type SendCustomerPendingOptions struct {
	OrderNumber string `json:"orderNumber"`
	StatusUrl   string `json:"statusUrl"`
}

func SendCustomerPending(ctx context.Context, order *db.Order) error {
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

	mailHTML, err := mail.Render(ctx, templateEmailPending, templateEmailPendingProps{
		OrderNumber: order.Number,
		TrackingURL: fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
	}, "e9042278-5556-4798-952e-34f1ce14dcf1", "47334277-ded2-4b14-a20c-a8a54557ae6b")
	if err != nil {
		return err
	}
	err = ses.Send(ctx, ses.SendProps{
		From:    "info@vietnam-immigrations.org",
		To:      lo.Compact([]string{"info@vietnam-immigrations.org", order.Billing.Email, order.Billing.Email2}),
		ReplyTo: "info@vietnam-immigrations.org",
		Subject: fmt.Sprintf("[Pending Review] Vietnam Visa Online Order #%s", order.Number),
		HTML:    *mailHTML,
	})
	return err
}
