package mail

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mail"
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

	mjmlUsername, err := ssm.GetParameter(ctx, "/mjml/username", false)
	if err != nil {
		return err
	}
	mjmlPassword, err := ssm.GetParameter(ctx, "/mjml/password", true)
	if err != nil {
		return err
	}
	mailHTML, err := mail.Render(ctx, templateEmailPriorityAdmin, templateEmailPriorityAdminProps{
		Entry:       order.Trip.Checkpoint,
		ArrivalDate: order.Trip.ArrivalDate,
		Applicants: lo.Map(order.Applicants, func(app db.Applicant, _ int) templateEmailPriorityAdminPropsApplicant {
			return templateEmailPriorityAdminPropsApplicant{
				RegistrationCode: app.RegistrationCode,
				Email:            app.Email,
			}
		}),
		ProcessingTime: processingTimeText,
		ExtraServices:  ftText,
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}

	err = ses.Send(ctx, ses.SendProps{
		From:        mailAddressInfo,
		To:          []string{cfg.EmailPartner},
		ReplyTo:     mailAddressInfo,
		CC:          []string{cfg.EmailPartnerCC},
		BCC:         nil,
		Subject:     subject,
		HTML:        *mailHTML,
		Attachments: nil,
	})

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
