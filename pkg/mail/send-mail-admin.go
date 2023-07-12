package mail

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/s3"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mail"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/amplify"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/notification"
)

func SendAdmin(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send mail to partner for order [%s]", order.Number)

	cfg, err := db.GetConfig(ctx)
	if err != nil {
		return err
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

	// extra services
	ft := make([]string, 0)
	if order.Trip.FastTrack != "No" {
		ft = append(ft, fmt.Sprintf("%s (flight: %s)", order.Trip.FastTrack, order.Trip.Flight))
	}
	if order.Trip.CarPickup {
		ft = append(ft, fmt.Sprintf("car pick-up (hotel: %s)", order.Trip.CarPickupAddress))
	}
	ftText := strings.Join(ft, ", ")

	mjmlUsername, err := secretsmanager.GetParameter(ctx, "/mjml/username")
	if err != nil {
		return err
	}
	mjmlPassword, err := secretsmanager.GetParameter(ctx, "/mjml/password")
	if err != nil {
		return err
	}
	mailHTML, err := mail.Render(ctx, templateEmailAdmin, templateEmailAdminProps{
		Entry:       order.Trip.Checkpoint,
		ArrivalDate: order.Trip.ArrivalDate,
		Applicants: lo.Map(order.Applicants, func(app db.Applicant, _ int) templateEmailAdminPropsApplicant {
			return templateEmailAdminPropsApplicant{
				LastName:  app.LastName,
				FirstName: app.FirstName,
			}
		}),
		ProcessingTime: processingTimeInContent,
		ExtraServices:  ftText,
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}

	err = ses.Send(ctx, ses.SendProps{
		From:    mailAddressInfo,
		To:      []string{cfg.EmailPartner},
		ReplyTo: mailAddressInfo,
		BCC:     nil,
		CC:      []string{cfg.EmailPartnerCC},
		Subject: subject,
		HTML:    *mailHTML,
		Attachments: lo.FlatMap(order.Applicants, func(app db.Applicant, _ int) []ses.SendPropsAttachment {
			portraitAtt, err := s3.ReadFileBucketSSM(ctx, amplify.S3Attachment, app.AttachmentPortrait.S3Key)
			if err != nil {
				log.Errorf("Failed to load portrait file [%s]: %s", app.AttachmentPortrait.S3Key, err)
				return nil
			}
			passportAtt, err := s3.ReadFileBucketSSM(ctx, amplify.S3Attachment, app.AttachmentPassport.S3Key)
			if err != nil {
				log.Errorf("Failed to load passport file [%s]: %s", app.AttachmentPassport.S3Key, err)
			}

			return []ses.SendPropsAttachment{
				{
					Name: fileNameFromS3Key(app.AttachmentPortrait.S3Key),
					Data: portraitAtt,
				},
				{
					Name: fileNameFromS3Key(app.AttachmentPassport.S3Key),
					Data: passportAtt,
				},
			}
		}),
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

func fileNameFromS3Key(key string) string {
	parts := strings.Split(key, "/")
	return parts[len(parts)-1]
}
