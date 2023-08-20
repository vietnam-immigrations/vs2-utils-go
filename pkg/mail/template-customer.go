package mail

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
}

const templateEmailCustomer = `<mjml>
    <mj-body background-color="#F4F4F4">
        <mj-spacer></mj-spacer>
        <mj-section background-color="#d4dfec">
            <mj-column>
                <mj-text
                        font-weight="bold"
                        font-size="20px"
                        align="center"
                        color="#55575d"
                >VIETNAM VISA ONLINE
                </mj-text>
            </mj-column>
        </mj-section>
        <mj-section background-color="white">
            <mj-column width="100%">
                <mj-text
                        color="#55575d"
                        font-weight="bold"
                        font-size="18px"
                >
                    Successfully Registered Application #{{.OrderNumber}}
                </mj-text>
                <mj-text
                        color="#55575d"
                        line-height="16px">
                    Thank you for choosing Vietnam-immigrations.org
                    <br/>
                    <br/>

                    We are pleased to confirm that your online visa application and associated payment have been successfully recorded. Below are the details of your application:
                </mj-text>
                <mj-spacer></mj-spacer>
                <mj-text
                    font-weight="bold" font-size="15px" color="#55575d"
                >
                    Order details
                </mj-text>
                <mj-table color="#55575d">
                    <tr>
                        <th>Arrival date</th>
                        <td>{{.ArrivalDate}}</td>
                    </tr>
                    <tr>
                        <th>Visa type</th>
                        <td>{{.VisaType}}</td>
                    </tr>
                    <tr>
                        <th>Visit purpose</th>
                        <td>{{.VisitPurpose}}</td>
                    </tr>
                    <tr>
                        <th>Entry</th>
                        <td>{{.Entry}}</td>
                    </tr>
                    <tr>
                        <th>Flight</th>
                        <td>{{.Flight}}</td>
                    </tr>
                    <tr>
                        <th>Hotel</th>
                        <td>{{.Hotel}}</td>
                    </tr>
                    <tr>
                        <th>Processing time</th>
                        <td>{{.ProcessingTime}}</td>
                    </tr>
                    <tr>
                        <th>Telephone</th>
                        <td>{{.Telephone}}</td>
                    </tr>
                    <tr>
                        <th>Email</th>
                        <td>{{.Email}}</td>
                    </tr>
                    <tr>
                        <th>Secondary email</th>
                        <td>{{.Email2}}</td>
                    </tr>
                </mj-table>
                {{if .HasExtraServices}}
                <mj-text
                        font-weight="bold" font-size="15px" color="#55575d"
                >
                    Extra services
                </mj-text>
                <mj-text color="#55575d">
                    {{.ExtraServices}}
                </mj-text>
                {{end}}
                <mj-text
                        font-weight="bold" font-size="15px" color="#55575d"
                >
                    Applicant(s)
                </mj-text>
                {{range .Applicants}}
                <mj-table color="#55575d">
                    <tr>
                        <th>Name</th><td>{{.Name}}</td>
                        <th>Nationality</th><td>{{.Nationality}}</td>
                    </tr>
                    <tr>
                        <th>Gender</th><td>{{.Gender}}</td>
                        <th>Passport</th><td>{{.Passport}}</td>
                    </tr>
                    <tr>
                        <th>Birthday</th><td>{{.Birthday}}</td>
                        <th>Valid until</th><td>{{.PassportValidUntil}}</td>
                    </tr>
                    <tr>
                        <th>Home address</th><td>{{.HomeAddress}}</td>
                        <th>Contact</th><td>{{.HomeContact}}</td>
                    </tr>
                    <tr>
                        <th>Address in Vietnam</th><td>{{.VietnamAddress}}</td>
                        <th>Previous visits</th><td>{{.PreviousVisitCount}}</td>
                    </tr>
                    <tr>
                        <th>Law violation</th><td>{{.LawViolation}}</td>
                    </tr>
                </mj-table>
                {{end}}
                <mj-spacer></mj-spacer>
                <mj-button
                        href="{{.TrackingURL}}"
                        background-color="#4d95ec"
                        color="white"
                >
                    Check your application status
                </mj-button>
                <mj-spacer></mj-spacer>
                <mj-divider border-width="1px" border-color="#E6E6E6"></mj-divider>
                <mj-spacer></mj-spacer>
            </mj-column>
            <mj-column width="100%">
                <mj-text
                        font-size="15px"
                        font-weight="bold"
                        color="#55575d"
                        align="center"
                >
                    Contact Us
                </mj-text>
            </mj-column>
            <mj-column width="50%">
                <mj-text
                        align="center"
                        font-weight="bold"
                        color="#55575d"
                >Email
                </mj-text>
                <mj-text align="center">
                    <a href="mailto:info@vietnam-immigrations.org">info@vietnam-immigrations.org</a>
                </mj-text>
            </mj-column>
            <mj-column width="50%">
                <mj-text
                        align="center"
                        font-weight="bold"
                        color="#55575d"
                >Telephone
                </mj-text>
                <mj-text
                        align="center"
                >
                    <a href="tel:+84888669951">(+84) 888 66 9951</a>
                    <br/>
                    <br/>
                    Also available in WhatsApp
                </mj-text>
            </mj-column>
        </mj-section>
        <mj-section
                background-color="#d4dfec"
        >
            <mj-column>
                <mj-text line-height="14px" font-size="12px" color="#55575d">
                    <ul>
                        <li>
                            The applicant(s) bears full responsibility for ensuring the accuracy of the passport and travel details provided above. Please note that amendments or cancellations are not allowed once the visa application is submitted in the immigration system.
                        </li>
                    </ul>
                </mj-text>
            </mj-column>
        </mj-section>
    </mj-body>
</mjml>`
