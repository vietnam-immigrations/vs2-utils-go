package mail

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/s3"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mail"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/amplify"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/db"
	"github.com/vietnam-immigrations/vs2-utils-go/v2/pkg/push"
)

func SendPriorityAdmin(ctx context.Context, order *db.Order) error {
	log := logger.FromContext(ctx)
	log.Infof("send priority email to partner for order [%s]", order.Number)
	enabled, err := checkSendMailAdminEnabled(ctx, order)
	if err != nil {
		log.Errorf("Failed to check send mail admin enabled: %s", err)
		return err
	}
	if enabled {
		return sendPriorityAdmin(ctx, order, nil)
	} else {
		log.Infof("Skip sending email to partner for order [%s]", order.Number)
		return sendPriorityAdmin(ctx, order, lo.ToPtr("info@vietnam-immigrations.org"))
	}
}

func sendPriorityAdmin(ctx context.Context, order *db.Order, overrideToEmail *string) error {
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
		VisaType:       order.Trip.VisaType,
		VisitPurpose:   order.Trip.VisitPurpose,
	}, mjmlUsername, mjmlPassword)
	if err != nil {
		return err
	}

	err = SendUseBrevo(ctx, ses.SendProps{
		From:    mailAddressInfo,
		To:      []string{lo.FromPtrOr(overrideToEmail, cfg.EmailPartner)},
		ReplyTo: mailAddressInfo,
		CC:      []string{cfg.EmailPartnerCC},
		BCC:     nil,
		Subject: subject,
		HTML:    *mailHTML,
		Attachments: lo.FlatMap(order.Applicants, func(app db.Applicant, _ int) []ses.SendPropsAttachment {
			attachments := make([]ses.SendPropsAttachment, 0)

			if app.AttachmentPortrait != nil && app.AttachmentPortrait.S3Key != "" {
				portraitAtt, err := s3.ReadFileBucketSSM(ctx, amplify.S3Attachment, app.AttachmentPortrait.S3Key)
				if err != nil {
					log.Errorf("Failed to load portrait file [%s]: %s", app.AttachmentPortrait.S3Key, err)
					return nil
				}
				attachments = append(attachments, ses.SendPropsAttachment{
					Name: fileNameFromS3Key(app.AttachmentPortrait.S3Key),
					Data: portraitAtt,
				})
			}

			if app.AttachmentPassport != nil && app.AttachmentPassport.S3Key != "" {
				passportAtt, err := s3.ReadFileBucketSSM(ctx, amplify.S3Attachment, app.AttachmentPassport.S3Key)
				if err != nil {
					log.Errorf("Failed to load passport file [%s]: %s", app.AttachmentPassport.S3Key, err)
				}
				attachments = append(attachments, ses.SendPropsAttachment{
					Name: fileNameFromS3Key(app.AttachmentPassport.S3Key),
					Data: passportAtt,
				})
			}

			return attachments
		}),
	})

	if err != nil {
		push.SendNotificationForOrder(
			ctx,
			order.ID.Hex(),
			"Send priority email to partner failed",
			fmt.Sprintf("Failed to send email to partner for order %s %s [#%s]: %s", order.Billing.LastName, order.Billing.FirstName, order.Number, err),
		)
		return err
	}
	push.SendNotificationForOrder(
		ctx,
		order.ID.Hex(),
		"Email sent to partner",
		fmt.Sprintf("Email sent to partner for order %s %s [#%s]", order.Billing.LastName, order.Billing.FirstName, order.Number),
	)
	return nil
}
