package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ses"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func SendUseBrevo(ctx context.Context, props ses.SendProps) error {
	log := logger.FromContext(ctx)

	brevoKey, err := ssm.GetParameter(ctx, "/brevo/key", true)
	if err != nil {
		return fmt.Errorf("failed to get brevo key: %w", err)
	}

	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", brevoKey)
	br := brevo.NewAPIClient(cfg)
	_, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{Email: props.From},
		To: lo.Map(props.To, func(email string, _ int) brevo.SendSmtpEmailTo {
			return brevo.SendSmtpEmailTo{Email: email}
		}),
		Bcc: lo.Map(props.BCC, func(email string, _ int) brevo.SendSmtpEmailBcc {
			return brevo.SendSmtpEmailBcc{Email: email}
		}),
		Cc: lo.Map(props.CC, func(email string, _ int) brevo.SendSmtpEmailCc {
			return brevo.SendSmtpEmailCc{Email: email}
		}),
		HtmlContent: props.HTML,
		Subject:     props.Subject,
		ReplyTo:     &brevo.SendSmtpEmailReplyTo{Email: props.ReplyTo},
		Attachment: lo.Map(props.Attachments, func(a ses.SendPropsAttachment, _ int) brevo.SendSmtpEmailAttachment {
			data := base64.StdEncoding.EncodeToString(a.Data)
			return brevo.SendSmtpEmailAttachment{
				Content: data,
				Name:    a.Name,
			}
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	if resp.StatusCode >= 400 {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Errorf("failed to close response body: %s", err)
			}
		}(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("status code %d, failed to send email, failed to read response body: %s", resp.StatusCode, err)
		}
		return fmt.Errorf("status code %d, failed to send email: %s", resp.StatusCode, string(body))
	}
	return nil
}
