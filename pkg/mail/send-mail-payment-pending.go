package mail

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mail"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	shopDB "github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/shop/db"
)

func SendPaymentPending(ctx context.Context, order *shopDB.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send email to customer for order [%s]", order.OrderNumber)

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
	mailHTML, err := mail.Render(ctx, templateEmailPaymentPending, templateEmailPaymentPendingProps{
		OrderNumber: order.OrderNumber,
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}
	err = ses.Send(ctx, ses.SendProps{
		From:    mailAddressInfo,
		To:      lo.Compact([]string{order.Billing.Email, order.Billing.Email2}),
		CC:      []string{cfg.EmailCustomerCC},
		ReplyTo: mailAddressInfo,
		Subject: fmt.Sprintf("[Payment Confirmation Pending] Vietnam Visa Online Order #%s", order.OrderNumber),
		HTML:    *mailHTML,
	})
	return err
}
