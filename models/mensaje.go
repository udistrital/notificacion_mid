package models

type Mensaje struct {
	Body          map[string]interface{}
	Attributes    map[string]string
	ReceiptHandle string
}
