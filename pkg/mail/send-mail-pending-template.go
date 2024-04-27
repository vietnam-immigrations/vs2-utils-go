package mail

import (
	_ "embed"
)

type templateEmailPendingProps struct {
	OrderNumber string
	TrackingURL string
}

//go:embed send-mail-pending-template.mjml
var templateEmailPending string
