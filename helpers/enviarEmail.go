package helpers

import (
	//"context"
	"encoding/json"
	//"strconv"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	//"time"
	//"github.com/udistrital/utils_oas/formatdata"
	//"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/udistrital/notificacion_mid/models"
)

func SendEmail(input models.SendEmailInput) (result *ses.SendRawEmailOutput, outputError map[string]interface{}) {

	// Attempt to send the email.
	if inputRaw, err := formatSendRawEmailInput(input); err == nil {

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)

		svc := ses.New(sess)

		//svc.CreateTemplate()
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

func SendTemplatedEmail(input models.InputTemplatedEmail) (result *ses.SendRawEmailOutput, outputError map[string]interface{}) {
	//formatdata.JsonPrint(input)
	if inputSesTemplate, errFormat := formatSendBulkTemplatedEmailInput(input); errFormat == nil {
		//formatdata.JsonPrint(inputSesTemplate)
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)

		svc := ses.New(sess)

		var testTemplate ses.TestRenderTemplateInput
		testTemplate.TemplateName = inputSesTemplate.Template

		for indexDestinations, dest := range inputSesTemplate.Destinations {
			//te := "{\"estado\":\"inscrito\",\"fecha\":\"2023-10-20\",\"nombre\":\"Fabian  David Barreto Sanchez\",\"periodo\":\"2024-1\"}"
			testTemplate.TemplateData = dest.ReplacementTemplateData
			if renderedTemplate, err := svc.TestRenderTemplate(&testTemplate); err == nil {
				mimeEmail := renderedTemplate.RenderedTemplate
				//fmt.Println(*mimeEmail)
				var ToAddressesString []string
				for _, address := range dest.Destination.ToAddresses {
					ToAddressesString = append(ToAddressesString, *address)
				}
				var rawEmail string
				rawEmail = ""
				rawEmail += "From:" + input.Source + "\n"
				rawEmail += "To:" + strings.Join(ToAddressesString, ",") + "\n"
				if len(input.Destinations[indexDestinations].Attachments) != 0 {
					rawEmail += addAttachmentsEmailMIME(*mimeEmail, input.Destinations[indexDestinations].Attachments)
				} else {

					rawEmail += *mimeEmail
				}
				rawMessage := ses.RawMessage{
					Data: []byte(rawEmail),
				}
				var inputRawEmail ses.SendRawEmailInput
				inputRawEmail.SetRawMessage(&rawMessage)
				resultado, err := svc.SendRawEmail(&inputRawEmail)
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
				}
			} else {
				fmt.Println("fail test", err)
				outputError = map[string]interface{}{"funcion": "/SendEmail/TestRenderTemplate", "err": err.Error(), "status": "502"}
			}
		}
	} else {
		outputError = errFormat
	}

	return
}

func formatSendBulkTemplatedEmailInput(input models.InputTemplatedEmail) (result ses.SendBulkTemplatedEmailInput, outputError map[string]interface{}) {

	result.Source = &input.Source
	result.Template = &input.Template

	for _, dest := range input.Destinations {
		var destSES ses.BulkEmailDestination
		destSES.Destination = dest.Destination
		if jsonReplaceData, err := json.Marshal(dest.ReplacementTemplateData); err == nil {
			templateDataString := strconv.Quote(string(jsonReplaceData))
			templateDataString = templateDataString[1 : len(templateDataString)-1]
			//fmt.Print(templateDataString)
			//fmt.Println(templateDataString)
			templateDataString = strings.ReplaceAll(templateDataString, `\`, ``)
			//fmt.Println("template data", templateDataString)
			destSES.ReplacementTemplateData = &templateDataString
			result.Destinations = append(result.Destinations, &destSES)
		} else {
			outputError = map[string]interface{}{"funcion": "/SendEmail/formatSendBulkTemplatedEmailInput", "err": err.Error(), "status": "502"}
		}
	}

	if jsonDefaultData, err := json.Marshal(input.DefaultTemplateData); err == nil {
		templateDataString := strconv.Quote(string(jsonDefaultData))
		templateDataString = templateDataString[1 : len(templateDataString)-1]
		templateDataString = strings.ReplaceAll(templateDataString, `\`, ``)
		result.DefaultTemplateData = &templateDataString
	} else {
		outputError = map[string]interface{}{"funcion": "/SendEmail/formatSendBulkTemplatedEmailInput", "err": err.Error(), "status": "502"}
	}
	return
}

func addAttachmentsEmailMIME(MIMEstring string, attachments []models.Attachment) (MIMEwithAttachments string) {
	pat := regexp.MustCompile(`boundary=".*."`)
	s := pat.FindString(MIMEstring)
	//fmt.Println(s)
	boundary := strings.Split(s, `"`)[1]
	indexfinalboundary := len(MIMEstring) - (len(boundary) + 6)
	var rawAttachments = ""
	if len(attachments) != 0 {
		for _, file := range attachments {
			rawAttachments += fmt.Sprintf("--%s\n", boundary)
			rawAttachments += fmt.Sprintf("Content-Type:"+*file.ContentType+"; name=\"%s\"\n", *file.FileName) // name
			rawAttachments += "Content-Transfer-Encoding:base64\n"
			rawAttachments += fmt.Sprintf("Content-Disposition:attachment;filename=\"%s\"\n", *file.FileName) // filename
			//rawAttachments += "Content-Disposition:inline;name={$file_name}\n"
			rawAttachments += fmt.Sprintf("Content-ID:<%s>\n", *file.FileName)
			rawAttachments += "\n"
			rawAttachments += *file.Base64File + "\n"
			rawAttachments += "\n"
		}
		MIMEwithAttachments = MIMEstring[:indexfinalboundary] + rawAttachments + MIMEstring[indexfinalboundary:]
	} else {
		MIMEwithAttachments = MIMEstring
	}
	return
}
