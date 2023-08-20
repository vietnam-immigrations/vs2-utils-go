package mail

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

const templateEmailAdmin = `<mjml>
    <mj-body background-color="#F4F4F4">
        <mj-spacer></mj-spacer>
        <mj-section background-color="white">
            <mj-column>
                <mj-text font-size="14px" line-height="20px">
                    {{.Entry}} on {{.ArrivalDate}}<br/>
                    {{.VisaType}} - {{.VisitPurpose}}<br/><br/><br/>
                    {{range .Applicants}}
                        <strong>{{.LastName}} {{.FirstName}}</strong><br/>
                        Home address: {{.AddressHome}}<br/>
                        Home contact: {{.ContactHome}}<br/>
                        Address in Vietnam: {{.AddressVietnam}}<br/>
                        Previous visit times: {{.PreviousVisitCount}}<br/>
                        Law violation: {{.LawViolation}}<br/><br/><br/>
                    {{end}}
                    {{.ProcessingTime}}<br/>
                    {{.ExtraServices}}
                </mj-text>
            </mj-column>
        </mj-section>
    </mj-body>
</mjml>`
