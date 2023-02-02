package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	mymailjet "github.com/nam-truong-le/lambda-utils-go/v3/pkg/mailjet"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/notification"
)

type SendPriorityAdminOptions struct {
	Applicants              []SendPriorityAdminOptionsApplicant `json:"applicants"`
	ArrivalDate             string                              `json:"arrivalDate"`
	Entry                   string                              `json:"entry"`
	ProcessingTimeInContent string                              `json:"processingTimeInContent"`
	ExtraServices           string                              `json:"extraServices"`
}

type SendPriorityAdminOptionsApplicant struct {
	RegistrationCode string `json:"registrationCode"`
	Email            string `json:"email"`
}

func SendPriorityAdmin(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send priority email to partner for order [%s]", order.Number)

	cfg, err := db.GetConfig(ctx)
	if err != nil {
		return err
	}

	to := mailjet.RecipientsV31{
		{
			Email: cfg.EmailPartner,
		},
	}

	applicants := make([]SendPriorityAdminOptionsApplicant, 0)
	for _, app := range order.Applicants {
		applicants = append(applicants, SendPriorityAdminOptionsApplicant{
			RegistrationCode: app.RegistrationCode,
			Email:            app.Email,
		})
	}

	// subject
	processingTimeText := processingTimeTexts[order.Trip.ProcessingTime]
	subject := fmt.Sprintf(
		"Khách duyệt web cục #%s %s %s %s [%s]",
		order.Number,
		order.Billing.LastName,
		order.Billing.FirstName,
		order.Trip.ArrivalDate,
		processingTimeText,
	)

	// extra services
	ft := make([]string, 0)
	if order.Trip.FastTrack != "No" {
		ft = append(ft, fmt.Sprintf("%s (flight: %s)", order.Trip.FastTrack, order.Trip.Flight))
	}
	ftText := strings.Join(ft, ", ")

	variables := SendPriorityAdminOptions{
		Applicants:              applicants,
		ArrivalDate:             order.Trip.ArrivalDate,
		Entry:                   order.Trip.Checkpoint,
		ProcessingTimeInContent: processingTimeText,
		ExtraServices:           ftText,
	}
	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		return err
	}
	rawVariables := new(map[string]interface{})
	err = json.Unmarshal(jsonVariables, rawVariables)
	if err != nil {
		return err
	}

	body := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: "info@vietnam-immigrations.org",
			Name:  "Vietnam Visa Online",
		},
		To: &to,
		Cc: &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: cfg.EmailPartnerCC,
			},
		},
		TemplateID:       cfg.PriorityEmailPartnerTemplateID,
		TemplateLanguage: true,
		Subject:          subject,
		Variables:        *rawVariables,
	}

	err = mymailjet.Send(ctx, body)
	if err != nil {
		_ = notification.Create(ctx, notification.Notification{
			ID:         uuid.New().String(),
			CSSClasses: "bg-negative text-white",
			Message:    fmt.Sprintf("Lỗi gửi email tới đối tác cho khách %s %s [#%s]", order.Billing.LastName, order.Billing.FirstName, order.Number),
		})
		return err
	}
	_ = notification.Create(ctx, notification.Notification{
		ID:         uuid.New().String(),
		CSSClasses: "bg-secondary text-white",
		Message:    fmt.Sprintf("Đã gửi email tới đối tác cho khách %s %s [#%s]", order.Billing.LastName, order.Billing.FirstName, order.Number),
	})
	return nil
}
