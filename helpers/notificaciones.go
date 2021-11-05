package helpers

import (
	"context"
	"fmt"
	"strconv"
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
	var listaDestinatarios string
	if len(body.DestinatarioId) > 1 {
		listaDestinatarios = "[" + strings.Join(body.DestinatarioId, ",") + "]"
	} else {
		listaDestinatarios = body.DestinatarioId[0]
	}
	logs.Debug(listaDestinatarios)

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

func Suscribir(body models.Suscripcion, atributos map[string]string) (Arn string, outputError map[string]interface{}) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/Suscribir", "err": err, "status": "502"}
		return "", outputError
	}

	client := sns.NewFromConfig(cfg)

	for _, subscriptor := range body.Suscritos {
		att := make(map[string]string)
		att["Destinatario"] = subscriptor.Id
		for k, v := range atributos {
			att[k] = v
		}
		input := &sns.SubscribeInput{
			Endpoint:              &subscriptor.Endpoint,
			Protocol:              aws.String(subscriptor.Protocolo),
			ReturnSubscriptionArn: true,
			TopicArn:              &body.ArnTopic,
			Attributes: map[string]string{
				"FilterPolicy": "{\"Destinatario\":[" + subscriptor.Id + "]}",
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

func ListaTopics() (topicArn []string, outputError map[string]interface{}) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err, "status": "502"}
		return nil, outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.ListTopicsInput{}

	results, err := client.ListTopics(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err, "status": "502"}
		return nil, outputError
	}

	for _, t := range results.Topics {
		topicArn = append(topicArn, *t.TopicArn)
	}
	return
}

func CrearTopic(nombre string, display string, fifo bool) (arn string, outputError map[string]interface{}) {
	if fifo {
		nombre += ".fifo"
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err, "status": "502"}
		return "", outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.CreateTopicInput{
		Name: &nombre,
		Attributes: map[string]string{
			"DisplayName": display,
			"FifoTopic":   strconv.FormatBool(fifo),
		},
	}

	results, err := client.CreateTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err, "status": "502"}
		return "", outputError
	}

	return *results.TopicArn, nil
}
