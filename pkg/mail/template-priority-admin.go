package mail

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
}

const templateEmailPriorityAdmin = `<mjml>
    <mj-body background-color="#F4F4F4">
        <mj-spacer></mj-spacer>
        <mj-section background-color="white">
            <mj-column>
                <mj-text font-size="14px" line-height="20px">
                    {{.Entry}} on {{.ArrivalDate}} <br/>
                    {{range .Applicants}}
                        {{.RegistrationCode}} | {{.Email}} <br/>
                    {{end}}
                    {{.ProcessingTime}} <br/>
                    {{.ExtraServices}}
                </mj-text>
            </mj-column>
        </mj-section>
    </mj-body>
</mjml>`
