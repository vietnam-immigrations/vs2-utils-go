package mail

import (
	_ "embed"
)

type templateEmailPriorityCustomerPropsApplicant struct {
	Title            string
	RegistrationCode string
}

type templateEmailPriorityCustomerProps struct {
	FullName    string
	ArrivalDate string
	Entry       string
	Applicants  []templateEmailPriorityCustomerPropsApplicant
}

//go:embed priority-send-mail-customer-template.mjml
var templateEmailPriorityCustomer string
