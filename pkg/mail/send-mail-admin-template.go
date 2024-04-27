package mail

import (
	_ "embed"
)

type templateEmailAdminPropsApplicant struct {
	LastName           string
	FirstName          string
	AddressHome        string
	ContactHome        string
	AddressVietnam     string
	PreviousVisitCount string
	LawViolation       string
}

type templateEmailAdminProps struct {
	VisaType       string
	VisitPurpose   string
	Entry          string
	ArrivalDate    string
	Applicants     []templateEmailAdminPropsApplicant
	ProcessingTime string
	ExtraServices  string
}

//go:embed send-mail-admin-template.mjml
var templateEmailAdmin string
