package mail

import (
	_ "embed"
)

type templateEmailImageRejectedProps struct {
	OrderNumber string
	UploadURL   string
}

//go:embed send-mail-images-rejected-template.mjml
var templateEmailImageRejected string
