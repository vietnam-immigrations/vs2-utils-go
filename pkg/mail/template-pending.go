package mail

type templateEmailPendingProps struct {
	OrderNumber string
	TrackingURL string
}

const templateEmailPending = `<mjml>
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
                    Application #{{.OrderNumber}}
                </mj-text>
                <mj-text
                        color="#55575d"
                        line-height="16px">
                    Thank you for choosing Vietnam-immigrations.org
                    <br/>
                    <br/>

                    We have successfully received your application and it is currently being reviewed along with the uploaded images. Once the review of your application data is complete, we will send you another confirmation email to notify you.
                </mj-text>
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
