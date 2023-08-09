package helpers

import (
	//"context"
	//"encoding/json"
	//"strconv"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	//"time"

	//"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/udistrital/notificacion_mid/models"
)

// const (
// 	// The subject line for the email.
// 	Subject = "Amazon SES Test (AWS SDK for Go)"

// 	// The HTML body for the email.
// 	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
// 		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
// 		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

// 	//The email body for recipients with non-HTML email clients.
// 	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

// 	// The character encoding for the email.
// 	CharSet = "UTF-8"
// )

func SendEmail(input models.SendEmailInput) (result *ses.SendRawEmailOutput, outputError map[string]interface{}) {

	// Attempt to send the email.
	if inputRaw, err := formatSendRawEmailInput(input); err == nil {

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)

		svc := ses.New(sess)

		resultado, err := svc.SendRawEmail(&inputRaw)

		result = resultado
		// Display error messages if they occur.
		if err != nil {
			outputError = map[string]interface{}{"funcion": "/SendEmail", "err": err.Error(), "status": "502"}
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ses.ErrCodeMessageRejected:
					fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
				case ses.ErrCodeMailFromDomainNotVerifiedException:
					fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
				case ses.ErrCodeConfigurationSetDoesNotExistException:
					fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
	} else {
		outputError = map[string]interface{}{"funcion": "/SendEmail/", "err": err["err"], "status": "502"}
	}

	return
}

func formatSendRawEmailInput(input models.SendEmailInput) (result ses.SendRawEmailInput, outputError map[string]interface{}) {
	hash := md5.Sum([]byte(*input.Message.Subject.Data))
	boundary := hex.EncodeToString(hash[:])

	var ToAddressesString []string
	for _, address := range input.Destination.ToAddresses {
		ToAddressesString = append(ToAddressesString, *address)
	}
	rawEmail := ""
	rawEmail += "From:" + *input.SourceName + "<" + *input.SourceEmail + ">\n"
	rawEmail += "To:" + strings.Join(ToAddressesString, ",") + "\n"
	rawEmail += fmt.Sprintf("Subject:%s\n", *input.Message.Subject.Data)
	rawEmail += "MIME-Version: 1.0\n"
	rawEmail += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n\n", boundary)
	rawEmail += fmt.Sprintf("--%s\n", boundary)

	rawEmail += fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"sub_%s\"\n\n", boundary)
	rawEmail += fmt.Sprintf("--sub_%s\n", boundary)
	rawEmail += "Content-Type: text/plain; charset=\"UTF-8\"\n"
	rawEmail += "\n"
	rawEmail += fmt.Sprintf("%s\n", *input.Message.Body.Text.Data)
	rawEmail += "\n"
	rawEmail += fmt.Sprintf("--sub_%s\n", boundary)
	rawEmail += "Content-Type: text/html; charset=\"UTF-8\"\n"
	rawEmail += "Content-Disposition: inline\n"
	rawEmail += "\n"
	rawEmail += fmt.Sprintf("%s\n", *input.Message.Body.Html.Data)
	rawEmail += "\n"
	rawEmail += fmt.Sprintf("--sub_%s--\n\n", boundary)
	if len(input.Message.Attachments) != 0 {
		for _, file := range input.Message.Attachments {
			rawEmail += fmt.Sprintf("--%s\n", boundary)
			rawEmail += fmt.Sprintf("Content-Type:"+*file.ContentType+"; name=\"%s\"\n", *file.FileName) // name
			rawEmail += "Content-Transfer-Encoding:base64\n"
			rawEmail += fmt.Sprintf("Content-Disposition:attachment;filename=\"%s\"\n", *file.FileName) // filename
			//rawEmail += "Content-Disposition:inline;name={$file_name}\n"
			rawEmail += fmt.Sprintf("Content-ID:<%s>\n", *file.FileName)
			rawEmail += "\n"
			rawEmail += *file.Base64File + "\n"
			rawEmail += "\n"
		}
	}
	rawEmail += fmt.Sprintf("--%s--\n\n", boundary)
	rawMessage := ses.RawMessage{
		Data: []byte(rawEmail),
	}

	result.SetRawMessage(&rawMessage)

	if err := result.Validate(); err != nil {

		outputError = map[string]interface{}{"funcion": "/SendEmail/formatSendRawEmailInput", "err": err.Error(), "status": "502"}
	}
	return
}
