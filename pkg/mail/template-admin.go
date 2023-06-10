package mail

type templateEmailAdminPropsApplicant struct {
	LastName  string
	FirstName string
}

type templateEmailAdminProps struct {
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
                    {{range .Applicants}}
                        {{.LastName}} {{.FirstName}}<br/>
                    {{end}}
                    {{.ProcessingTime}}<br/>
                    {{.ExtraServices}}
                </mj-text>
            </mj-column>
        </mj-section>
    </mj-body>
</mjml>`
