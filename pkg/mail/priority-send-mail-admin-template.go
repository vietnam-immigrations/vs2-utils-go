package mail

import (
	_ "embed"
)

type templateEmailPriorityAdminPropsApplicant struct {
	RegistrationCode string
	Email            string
}

type templateEmailPriorityAdminProps struct {
	Entry          string
	ArrivalDate    string
	Applicants     []templateEmailPriorityAdminPropsApplicant
	ProcessingTime string
	ExtraServices  string
	VisaType       string
	VisitPurpose   string
}

//go:embed priority-send-mail-admin-template.mjml
var templateEmailPriorityAdmin string
