package mail

import (
	_ "embed"
)

type templateEmailPaymentPendingProps struct {
	OrderNumber string
}

//go:embed send-mail-payment-pending-template.mjml
var templateEmailPaymentPending string
