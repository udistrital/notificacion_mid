package models

import (
	"github.com/aws/aws-sdk-go/service/ses"
)

type InputTemplatedEmail struct {
	Source              string
	Template            string
	Destinations        []DestinationTemplate
	DefaultTemplateData map[string]interface{}
}

type DestinationTemplate struct {
	Destination             *ses.Destination
	ReplacementTemplateData map[string]interface{}
	Attachments             []Attachment
}
