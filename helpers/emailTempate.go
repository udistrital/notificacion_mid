package helpers

import (
	//"context"
	//"encoding/json"
	//"strconv"
	//"crypto/md5"
	//"encoding/hex"
	"fmt"
	//"strconv"
	//"strings"

	//"time"
	//"github.com/udistrital/utils_oas/formatdata"
	//"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	//"github.com/udistrital/notificacion_mid/models"
)

func CreateEmailTemplate(input ses.CreateTemplateInput) (result *ses.CreateTemplateOutput, outputError map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := ses.New(sess)

	res, err := svc.CreateTemplate(&input)

	if err != nil {
		outputError = map[string]interface{}{"funcion": "/CreateEmailTemplate", "err": err.Error(), "status": "502"}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeAlreadyExistsException:
				fmt.Println(ses.ErrCodeAlreadyExistsException, aerr.Error())
			case ses.ErrCodeInvalidTemplateException:
				fmt.Println(ses.ErrCodeInvalidTemplateException, aerr.Error())
			case ses.ErrCodeLimitExceededException:
				fmt.Println(ses.ErrCodeLimitExceededException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, outputError
	} else {
		result = res
	}
	return
}

func GetEmailTemplate(input ses.GetTemplateInput) (result *ses.GetTemplateOutput, outputError map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		outputError = map[string]interface{}{"funcion": "/GetEmailTemplate", "err": err.Error(), "status": "502"}
	} else {
		svc := ses.New(sess)

		if result, err := svc.GetTemplate(&input); err == nil {
			return result, nil
		} else {
			outputError = map[string]interface{}{"funcion": "/GetEmailTemplate", "err": err.Error(), "status": "502"}
			return nil, outputError
		}
	}
	return
}

func ListEmailTemplates(input ses.ListTemplatesInput) (result *ses.ListTemplatesOutput, outputError map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		outputError = map[string]interface{}{"funcion": "/GetEmailTemplate", "err": err.Error(), "status": "502"}

	} else {
		svc := ses.New(sess)

		if result, err := svc.ListTemplates(&input); err == nil {
			return result, nil
		} else {
			outputError = map[string]interface{}{"funcion": "/GetEmailTemplate", "err": err.Error(), "status": "502"}
			return nil, outputError
		}
	}
	return
}

func UpdateEmailTemplate(input ses.UpdateTemplateInput) (result *ses.UpdateTemplateOutput, outputError map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		outputError = map[string]interface{}{"funcion": "/UpdateEmailTemplate", "err": err.Error(), "status": "502"}
	} else {
		svc := ses.New(sess)

		if result, err := svc.UpdateTemplate(&input); err == nil {
			return result, nil
		} else {
			outputError = map[string]interface{}{"funcion": "/UpdateEmailTemplate", "err": err.Error(), "status": "502"}
			return nil, outputError
		}
	}
	return
}

func DeleteEmailTemplate(input ses.DeleteTemplateInput) (result *ses.DeleteTemplateOutput, outputError map[string]interface{}) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		outputError = map[string]interface{}{"funcion": "/GetEmailTemplate", "err": err.Error(), "status": "502"}
	} else {
		svc := ses.New(sess)

		if result, err := svc.DeleteTemplate(&input); err == nil {
			return result, nil
		} else {
			outputError = map[string]interface{}{"funcion": "/GetEmailTemplate", "err": err.Error(), "status": "502"}
			return nil, outputError
		}
	}
	return
}
