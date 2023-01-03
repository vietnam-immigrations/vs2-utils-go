package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
	mymailjet "github.com/nam-truong-le/lambda-utils-go/v2/pkg/mailjet"
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

	applicants := make([]SendPriorityCustomerOptionsApplicant, 0)
	for i, app := range order.Applicants {
		applicants = append(applicants, SendPriorityCustomerOptionsApplicant{
			Title:            fmt.Sprintf("Applicant %d", i+1),
			RegistrationCode: app.RegistrationCode,
		})
	}

	// extra services
	ft := make([]string, 0)
	if order.Trip.FastTrack != "No" {
		ft = append(ft, fmt.Sprintf("%s (flight: %s)", order.Trip.FastTrack, order.Trip.Flight))
	}
	ftText := strings.Join(ft, ", ")

	variables := SendPriorityCustomerOptions{
		FullName:      strings.ToUpper(fmt.Sprintf("%s %s", order.Billing.FirstName, order.Billing.LastName)),
		Applicants:    applicants,
		ExtraServices: ftText,
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
		TemplateID:       cfg.PriorityEmailCustomerTemplateID,
		TemplateLanguage: true,
		Subject:          fmt.Sprintf("Vietnam Visa Priority Handling #%s", order.Number),
		Variables:        *rawVariables,
	}

	return mymailjet.Send(ctx, body)
}
