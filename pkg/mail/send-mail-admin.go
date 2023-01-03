package mail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mailjet/mailjet-apiv3-go/v4"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/aws/s3"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
	mymailjet "github.com/nam-truong-le/lambda-utils-go/v2/pkg/mailjet"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/aws/ssm"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/notification"
)

type SendAdminOptionsApplicant struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type SendAdminOptions struct {
	Applicants              []SendAdminOptionsApplicant `json:"applicants"`
	ArrivalDate             string                      `json:"arrivalDate"`
	Entry                   string                      `json:"entry"`
	ProcessingTimeInContent string                      `json:"processingTimeInContent"`
	ExtraServices           string                      `json:"extraServices"`
}

func SendAdmin(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send mail to partner for order [%s]", order.Number)

	cnf, err := db.GetConfig(ctx)
	if err != nil {
		return err
	}

	to := mailjet.RecipientsV31{
		{
			Email: cnf.EmailPartner,
		},
	}

	applicants := make([]SendAdminOptionsApplicant, 0)
	for _, app := range order.Applicants {
		applicants = append(applicants, SendAdminOptionsApplicant{
			FirstName: app.FirstName,
			LastName:  app.LastName,
		})
	}

	// subject
	subject := fmt.Sprintf(
		"#%s %s %s %s [%s]",
		order.Number,
		order.Billing.LastName,
		order.Billing.FirstName,
		order.Trip.ArrivalDate,
		processingTimeTexts[order.Trip.ProcessingTime],
	)
	processingTimeInContent := processingTimeTexts[order.Trip.ProcessingTime]

	// attachments
	attachments := make(mailjet.AttachmentsV31, 0)
	for _, app := range order.Applicants {
		portraitAtt, err := s3.ReadFileBucketSSM(ctx, ssm.S3BucketAttachment, app.AttachmentPortrait.S3Key)
		if err != nil {
			return err
		}
		attachments = append(attachments, mailjet.AttachmentV31{
			Base64Content: base64.StdEncoding.EncodeToString(portraitAtt),
			Filename:      fileNameFromS3Key(app.AttachmentPortrait.S3Key),
			ContentType:   "application/octet-stream",
		})

		passportAtt, err := s3.ReadFileBucketSSM(ctx, ssm.S3BucketAttachment, app.AttachmentPassport.S3Key)
		attachments = append(attachments, mailjet.AttachmentV31{
			Base64Content: base64.StdEncoding.EncodeToString(passportAtt),
			Filename:      fileNameFromS3Key(app.AttachmentPassport.S3Key),
			ContentType:   "application/octet-stream",
		})
	}

	// extra services
	ft := make([]string, 0)
	if order.Trip.FastTrack != "No" {
		ft = append(ft, fmt.Sprintf("%s (flight: %s)", order.Trip.FastTrack, order.Trip.Flight))
	}
	if order.Trip.CarPickup {
		ft = append(ft, fmt.Sprintf("car pick-up (hotel: %s)", order.Trip.CarPickupAddress))
	}
	ftText := strings.Join(ft, ", ")

	variables := SendAdminOptions{
		Applicants:              applicants,
		ArrivalDate:             order.Trip.ArrivalDate,
		Entry:                   order.Trip.Checkpoint,
		ProcessingTimeInContent: processingTimeInContent,
		ExtraServices:           ftText,
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

	body := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: "info@vietnam-immigrations.org",
			Name:  "Vietnam Visa Online",
		},
		To: &to,
		Cc: &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: cnf.EmailPartnerCC,
			},
		},
		TemplateID:       cnf.EmailPartnerTemplateID,
		TemplateLanguage: true,
		Subject:          subject,
		Attachments:      &attachments,
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

func fileNameFromS3Key(key string) string {
	parts := strings.Split(key, "/")
	return parts[len(parts)-1]
}
