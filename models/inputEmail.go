package models

import (
	"github.com/aws/aws-sdk-go/service/ses"
)

type SendEmailInput struct {
	_ struct{} `type:"structure"`

	// The name of the configuration set to use when you send an email using SendEmail.
	ConfigurationSetName *string `type:"string"`

	// The destination for this email, composed of To:, CC:, and BCC: fields.
	//
	// Destination is a required field
	Destination *Destination `type:"structure" required:"true"`

	// The message to be sent.
	//
	// Message is a required field
	Message *Message `type:"structure" required:"true"`

	// The reply-to email address(es) for the message. If the recipient replies
	// to the message, each reply-to address will receive the reply.
	ReplyToAddresses []*string `type:"list"`

	// The email address that bounces and complaints will be forwarded to when feedback
	// forwarding is enabled. If the message cannot be delivered to the recipient,
	// then an error message will be returned from the recipient's ISP; this message
	// will then be forwarded to the email address specified by the ReturnPath parameter.
	// The ReturnPath parameter is never overwritten. This email address must be
	// either individually verified with Amazon SES, or from a domain that has been
	// verified with Amazon SES.
	ReturnPath *string `type:"string"`

	// This parameter is used only for sending authorization. It is the ARN of the
	// identity that is associated with the sending authorization policy that permits
	// you to use the email address specified in the ReturnPath parameter.
	//
	// For example, if the owner of example.com (which has ARN arn:aws:ses:us-east-1:123456789012:identity/example.com)
	// attaches a policy to it that authorizes you to use feedback@example.com,
	// then you would specify the ReturnPathArn to be arn:aws:ses:us-east-1:123456789012:identity/example.com,
	// and the ReturnPath to be feedback@example.com.
	//
	// For more information about sending authorization, see the Amazon SES Developer
	// Guide (https://docs.aws.amazon.com/ses/latest/DeveloperGuide/sending-authorization.html).
	ReturnPathArn *string `type:"string"`

	// The email address that is sending the email. This email address must be either
	// individually verified with Amazon SES, or from a domain that has been verified
	// with Amazon SES. For information about verifying identities, see the Amazon
	// SES Developer Guide (https://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-addresses-and-domains.html).
	//
	// If you are sending on behalf of another user and have been permitted to do
	// so by a sending authorization policy, then you must also specify the SourceArn
	// parameter. For more information about sending authorization, see the Amazon
	// SES Developer Guide (https://docs.aws.amazon.com/ses/latest/DeveloperGuide/sending-authorization.html).
	//
	// Amazon SES does not support the SMTPUTF8 extension, as described in RFC6531
	// (https://tools.ietf.org/html/rfc6531). For this reason, the local part of
	// a source email address (the part of the email address that precedes the @
	// sign) may only contain 7-bit ASCII characters (https://en.wikipedia.org/wiki/Email_address#Local-part).
	// If the domain part of an address (the part after the @ sign) contains non-ASCII
	// characters, they must be encoded using Punycode, as described in RFC3492
	// (https://tools.ietf.org/html/rfc3492.html). The sender name (also known as
	// the friendly name) may contain non-ASCII characters. These characters must
	// be encoded using MIME encoded-word syntax, as described in RFC 2047 (https://tools.ietf.org/html/rfc2047).
	// MIME encoded-word syntax uses the following form: =?charset?encoding?encoded-text?=.
	//
	// SourceName is a required field
	SourceName *string `type:"string" required:"true"`

	// Source email origin
	//
	// SourceEmail is a required field
	SourceEmail *string `type:"string" required:"true"`

	// This parameter is used only for sending authorization. It is the ARN of the
	// identity that is associated with the sending authorization policy that permits
	// you to send for the email address specified in the Source parameter.
	//
	// For example, if the owner of example.com (which has ARN arn:aws:ses:us-east-1:123456789012:identity/example.com)
	// attaches a policy to it that authorizes you to send from user@example.com,
	// then you would specify the SourceArn to be arn:aws:ses:us-east-1:123456789012:identity/example.com,
	// and the Source to be user@example.com.
	//
	// For more information about sending authorization, see the Amazon SES Developer
	// Guide (https://docs.aws.amazon.com/ses/latest/DeveloperGuide/sending-authorization.html).
	SourceArn *string `type:"string"`

	// A list of tags, in the form of name/value pairs, to apply to an email that
	// you send using SendEmail. Tags correspond to characteristics of the email
	// that you define, so that you can publish email sending events.
	Tags []*ses.MessageTag `type:"list"`
}

type Destination struct {
	_ struct{} `type:"structure"`

	// The recipients to place on the BCC: line of the message.
	BccAddresses []*string `type:"list"`

	// The recipients to place on the CC: line of the message.
	CcAddresses []*string `type:"list"`

	// The recipients to place on the To: line of the message.
	ToAddresses []*string `type:"list"`
}

type Message struct {
	_ struct{} `type:"structure"`

	// The message body.
	//
	// Body is a required field
	Body *ses.Body `type:"structure" required:"true"`

	// The subject of the message: A short summary of the content, which will appear
	// in the recipient's inbox.
	//
	// Subject is a required field
	Subject *ses.Content `type:"structure" required:"true"`

	// The attachments include files
	Attachments []*Attachment `type:"list"`
}

type Attachment struct {
	ContentType *string
	FileName    *string
	Base64File  *string
}
