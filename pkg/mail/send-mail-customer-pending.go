package mail

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
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

	mjmlUsername, err := ssm.GetParameter(ctx, "/mjml/username", false)
	if err != nil {
		return err
	}
	mjmlPassword, err := ssm.GetParameter(ctx, "/mjml/password", true)
	if err != nil {
		return err
	}
	mailHTML, err := mail.Render(ctx, templateEmailPending, templateEmailPendingProps{
		OrderNumber: order.Number,
		TrackingURL: fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
	}, mjmlUsername, mjmlPassword)
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
