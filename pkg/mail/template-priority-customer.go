package mail

type templateEmailPriorityCustomerPropsApplicant struct {
	Title            string
	RegistrationCode string
}

type templateEmailPriorityCustomerProps struct {
	FullName   string
	Applicants []templateEmailPriorityCustomerPropsApplicant
}

const templateEmailPriorityCustomer = `<mjml>
    <mj-body background-color="#F4F4F4">
        <mj-spacer></mj-spacer>
        <mj-section background-color="white">
            <mj-column>
                <mj-text font-size="14px" line-height="20px">
                    <p>Dear {{.FullName}}</p>
                    <p>Thank you for your request to prioritize your visa application on the Immigration Portal.
                        We are pleased to inform you that your payment has been successfully processed, and we
                        have recorded the following details:</p>
                    <p style="color: darkblue">Information details:</p>
                    {{range .Applicants}}
                        <p>
                            <strong>{{.Title}}</strong>
                        </p>
                        <ul>
                            <li>Registration code: {{.RegistrationCode}}</li>
                        </ul>
                    {{end}}
                    <p>
                        <strong>Important:</strong> Please note that as the applicants, you are solely
                        responsible for the accuracy of the information provided during the registration
                        process on the Immigration Portal. Our system expedites the E-visa issuance
                        process, but it cannot identify or rectify any errors made by the applicants.
                        Once the E-visa is initiated or officially issued to the applicant, it cannot
                        be amended or cancelled for a refund, as per the Immigration Rules. In the event
                        that your application is denied by the Immigration Department, no E-visa or
                        refund will be issued.
                    </p>
                    <p>
                        Please be aware that the E-visa process is available during our business hours
                        from 8:30 AM to 3:30 PM (UTC+7), Monday to Friday.
                    </p>
                    <p><strong>Once the E-visa is granted, it cannot be modified.</strong></p>
                    <p><strong>Best Regards,</strong></p>
                    <p><strong>
                        E-visa Support Center<br/>
                        Vietnam Immigration Office<br/>
                        Hotline: (+84) 888 66 99 51<br/>
                        Email: info@vietnam-immigrations.org
                    </strong></p>
                </mj-text>
            </mj-column>
        </mj-section>
    </mj-body>
</mjml>`
