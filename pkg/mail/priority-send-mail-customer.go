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

type SendPriorityCustomerOptionsApplicant struct {
	Title            string `json:"title"`
	RegistrationCode string `json:"registrationCode"`
}

type SendPriorityCustomerOptions struct {
	FullName      string                                 `json:"fullName"`
	Applicants    []SendPriorityCustomerOptionsApplicant `json:"applicants"`
	ExtraServices string                                 `json:"extraServices"`
}

func SendPriorityCustomer(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send priority email to customer for order [%s]", order.Number)

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
	mailHTML, err := mail.Render(ctx, templateEmailPriorityCustomer, templateEmailPriorityCustomerProps{
		FullName: strings.ToUpper(fmt.Sprintf("%s %s", order.Billing.FirstName, order.Billing.LastName)),
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
	err = ses.Send(ctx, ses.SendProps{
		From:        mailAddressInfo,
		To:          lo.Compact([]string{order.Billing.Email, order.Billing.Email2}),
		ReplyTo:     mailAddressInfo,
		BCC:         nil,
		CC:          []string{cfg.EmailCustomerCC},
		Subject:     fmt.Sprintf("Vietnam Visa Priority Handling #%s", order.Number),
		HTML:        *mailHTML,
		Attachments: nil,
	})

	return err
}
