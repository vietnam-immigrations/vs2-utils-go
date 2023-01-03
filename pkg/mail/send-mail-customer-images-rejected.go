package mail

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
	mymailjet "github.com/nam-truong-le/lambda-utils-go/v2/pkg/mailjet"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
)

type SendCustomerImagesRejectedOptions struct {
	OrderNumber string `json:"orderNumber"`
	StatusUrl   string `json:"statusUrl"`
}

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

	variables := SendCustomerImagesRejectedOptions{
		OrderNumber: order.Number,
		StatusUrl:   fmt.Sprintf("https://%s/#/?order=%s&secret=%s", cfg.CustomerDomain, order.Number, order.OrderKey),
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

	log.Infof("%+v", rawVariables)

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
		TemplateID:       cfg.EmailCustomerRejectedImagesTemplateID,
		TemplateLanguage: true,
		Subject:          fmt.Sprintf("[IMPORTANT - PLEASE PROVIDE NEW IMAGES] Vietnam Visa Online Order #%s", order.Number),
		Variables:        *rawVariables,
	}

	return mymailjet.Send(ctx, body)
}
