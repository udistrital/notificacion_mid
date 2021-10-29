package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/udistrital/notificacion_api/models"
)

func PublicarNotificacion(body models.Notificacion) (msgId string, outputError map[string]interface{}) {
	tipoString := "String"
	tipoLista := "String.Array"
	listaDestinatarios := strings.Join(body.DestinatarioId, ",")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err, "status": "502"}
		return "", outputError
	}
	client := sns.NewFromConfig(cfg)

	atributos := make(map[string]types.MessageAttributeValue)
	atributos["Remitente"] = types.MessageAttributeValue{
		DataType:    &tipoString,
		StringValue: &body.RemitenteId,
	}
	atributos["Destinatario"] = types.MessageAttributeValue{
		DataType:    &tipoLista,
		StringValue: &listaDestinatarios,
	}

	for key, value := range body.Atributos {
		str := fmt.Sprintf("%v", value)
		atributos[key] = types.MessageAttributeValue{
			DataType:    &tipoString,
			StringValue: &str,
		}
	}

	input := &sns.PublishInput{
		Message:           &body.Mensaje,
		MessageAttributes: atributos,
		Subject:           &body.Asunto,
		TopicArn:          &body.Arn,
	}

	result, err := client.Publish(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err, "status": "502"}
		return "", outputError
	}

	msgId = *result.MessageId

	return
}

func Suscribir(body models.Suscripcion) (Arn string, outputError map[string]interface{}) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/Suscribir", "err": err, "status": "502"}
		return "", outputError
	}

	client := sns.NewFromConfig(cfg)

	for _, subscriptor := range body.Suscritos {
		input := &sns.SubscribeInput{
			Endpoint:              &subscriptor.Endpoint,
			Protocol:              aws.String(subscriptor.Protocolo),
			ReturnSubscriptionArn: true,
			TopicArn:              &body.ArnTopic,
			Attributes: map[string]string{
				"Destinatario": subscriptor.Id,
			},
		}

		result, err := client.Subscribe(context.TODO(), input)
		if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err, "status": "502"}
			return "", outputError
		}
		Arn = *result.SubscriptionArn
	}

	return
}
