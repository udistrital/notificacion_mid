package helpers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/udistrital/notificacion_mid/models"
)

func PublicarNotificacion(body models.Notificacion) (msgId string, outputError map[string]interface{}) {
	tipoString := "String"
	tipoLista := "String.Array"
	var listaDestinatarios string
	if len(body.DestinatarioId) > 1 {
		listaDestinatarios = "[\"" + strings.Join(body.DestinatarioId, "\",\"") + "\"]"
	} else {
		listaDestinatarios = "\"" + body.DestinatarioId[0] + "\""
	}

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

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

	var input *sns.PublishInput

	if strings.Contains(body.Arn, ".fifo") {
		input = &sns.PublishInput{
			Message:                &body.Mensaje,
			MessageAttributes:      atributos,
			Subject:                &body.Asunto,
			TopicArn:               &body.Arn,
			MessageDeduplicationId: &body.IdDeduplicacion,
			MessageGroupId:         &body.IdGrupoMensaje,
		}
	} else {
		input = &sns.PublishInput{
			Message:           &body.Mensaje,
			MessageAttributes: atributos,
			Subject:           &body.Asunto,
			TopicArn:          &body.Arn,
		}
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

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Suscribir", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

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
				"FilterPolicy": "{\"Destinatario\":[\"" + subscriptor.Id + "\",\"todos\"]}",
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

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

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

func CrearTopic(topic models.Topic) (arn string, outputError map[string]interface{}) {
	var tags []types.Tag
	var key0 string = "Name"
	var key1 string = "Environment"
	var val1 string = "prod"

	if beego.BConfig.RunMode == "dev" || beego.BConfig.RunMode == "test" {
		val1 = "test"
	}

	if topic.Fifo {
		topic.Nombre += ".fifo"
	}

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err, "status": "502"}
		return "", outputError
	}

	tags = append(tags, types.Tag{
		Key:   &key0,
		Value: &topic.Nombre,
	}, types.Tag{
		Key:   &key1,
		Value: &val1,
	})

	client := sns.NewFromConfig(cfg)

	input := &sns.CreateTopicInput{
		Name: &topic.Nombre,
		Attributes: map[string]string{
			"DisplayName": topic.Display,
			"FifoTopic":   strconv.FormatBool(topic.Fifo),
		},
		Tags: tags,
	}

	results, err := client.CreateTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err, "status": "502"}
		return "", outputError
	}

	return *results.TopicArn, nil
}

func VerificarSuscripcion(arn string, endpoint string) (suscrito bool, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/VerificarSuscripcion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err, "status": "502"}
		return false, outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.ListSubscriptionsByTopicInput{
		TopicArn: &arn,
	}

	results, err := client.ListSubscriptionsByTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err, "status": "502"}
		return false, outputError
	}

	for _, resultado := range results.Subscriptions {
		if *resultado.Endpoint == endpoint {
			return true, nil
		}
	}
	return
}

func BorrarTopic(arn string) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BorrarTopic", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarTopic", "err": err, "status": "502"}
		return outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.DeleteTopicInput{
		TopicArn: &arn,
	}

	_, err = client.DeleteTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarTopic", "err": err, "status": "502"}
		return outputError
	}

	return
}
