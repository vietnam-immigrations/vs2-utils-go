package mail

import (
	_ "embed"
)

type templateEmailCustomerPropsApplicant struct {
	Name               string
	Gender             string
	Nationality        string
	Passport           string
	Birthday           string
	PassportValidUntil string
	HomeAddress        string
	HomeContact        string
	VietnamAddress     string
	PreviousVisitCount string
	LawViolation       string
}

type templateEmailCustomerProps struct {
	OrderNumber      string
	VisaType         string
	VisitPurpose     string
	ArrivalDate      string
	Entry            string
	Flight           string
	Hotel            string
	ProcessingTime   string
	Telephone        string
	Email            string
	Email2           string
	HasExtraServices bool
	ExtraServices    string
	Applicants       []templateEmailCustomerPropsApplicant
	TrackingURL      string
	IsVisaOnArrival  bool
}

//go:embed send-mail-customer-template.mjml
var templateEmailCustomer string
