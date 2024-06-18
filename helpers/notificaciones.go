package helpers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/udistrital/notificacion_mid/models"
)

func PublicarLote(loteMensajes []models.Mensaje) (outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PublicarLote", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	// Cargar configuración
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PublicarLote", "err": err.Error(), "status": "502"}
		return outputError
	}
	client := sns.NewFromConfig(cfg)

	// Crear entradas de publicación
	entries := make([]types.PublishBatchRequestEntry, len(loteMensajes))
	for i, msg := range loteMensajes {
		body := msg.Body
		message, subject := body["Message"].(string), body["Subject"].(string)
		messageAttributes := body["MessageAttributes"].(map[string]interface{})
		usuarioDestino := messageAttributes["UsuarioDestino"].(map[string]interface{})["Value"].(string)
		idDuplicacion := fmt.Sprintf("%d", time.Now().UnixNano())

		// Construir atributos del mensaje
		atributos := make(map[string]types.MessageAttributeValue)
		for key, value := range messageAttributes {
			dato := value.(map[string]interface{})
			tipoDato := dato["Type"].(string)
			valorDato := dato["Value"].(string)
			atributos[key] = types.MessageAttributeValue{
				DataType:    &tipoDato,
				StringValue: &valorDato,
			}
		}

		// Entradas de publicación (lista de mensajes)
		entries[i] = types.PublishBatchRequestEntry{
			Id:                     aws.String(fmt.Sprintf("%v", i)),
			Message:                aws.String(message),
			MessageDeduplicationId: aws.String(idDuplicacion),
			MessageGroupId:         aws.String(usuarioDestino),
			Subject:                aws.String(subject),
			MessageAttributes:      atributos,
		}
	}

	// Publicar lote de mensajes
	_, err = client.PublishBatch(context.TODO(), &sns.PublishBatchInput{
		PublishBatchRequestEntries: entries,
		TopicArn:                   aws.String(loteMensajes[0].Body["TopicArn"].(string)),
	})
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PublicarLote", "err": err.Error(), "status": "502"}
		return outputError
	}
	return
}

func Publicar(body models.Notificacion, retornarInput bool) (msgId interface{}, outputError map[string]interface{}) {
	tipoString := "String"
	tipoLista := "String.Array"
	atributos := make(map[string]types.MessageAttributeValue)
	var listaDestinatarios string
	if len(body.DestinatarioId) > 1 {
		listaDestinatarios = "[\"" + strings.Join(body.DestinatarioId, "\",\"") + "\"]"
		atributos["Destinatario"] = types.MessageAttributeValue{
			DataType:    &tipoLista,
			StringValue: &listaDestinatarios,
		}
	} else if len(body.DestinatarioId) == 1 {
		listaDestinatarios = "\"" + body.DestinatarioId[0] + "\""
		atributos["Destinatario"] = types.MessageAttributeValue{
			DataType:    &tipoLista,
			StringValue: &listaDestinatarios,
		}
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
		outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err.Error(), "status": "502"}
		return "", outputError
	}
	client := sns.NewFromConfig(cfg)

	atributos["Remitente"] = types.MessageAttributeValue{
		DataType:    &tipoString,
		StringValue: &body.RemitenteId,
	}

	for key, value := range body.Atributos {
		str := fmt.Sprintf("%v", value)
		atributos[key] = types.MessageAttributeValue{
			DataType:    &tipoString,
			StringValue: &str,
		}
	}
	logs.Debug(body.Atributos)
	logs.Debug(&atributos)

	var input *sns.PublishInput

	if strings.Contains(body.ArnTopic, ".fifo") {
		input = &sns.PublishInput{
			Message:                &body.Mensaje,
			MessageAttributes:      atributos,
			Subject:                &body.Asunto,
			TopicArn:               &body.ArnTopic,
			MessageDeduplicationId: &body.IdDeduplicacion,
			MessageGroupId:         &body.IdGrupoMensaje,
		}
	} else {
		input = &sns.PublishInput{
			Message:           &body.Mensaje,
			MessageAttributes: atributos,
			Subject:           &body.Asunto,
			TopicArn:          &body.ArnTopic,
		}
	}
	result, err := client.Publish(context.TODO(), input, func(o *sns.Options) {
		o.ClientLogMode = aws.LogRequestWithBody | aws.LogResponseWithBody
	})
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err.Error(), "status": "502"}
		return "", outputError
	}

	msgId = *result.MessageId
	if retornarInput {
		inputMap := map[string]interface{}{
			"Message":           body.Mensaje,
			"MessageAttributes": getAtributosConsulta(atributos),
			"MessageId":         *result.MessageId,
			"Subject":           body.Asunto,
			"TopicArn":          body.ArnTopic,
		}
		msgId = inputMap
	}

	return
}

func getAtributosConsulta(atributos map[string]types.MessageAttributeValue) map[string]interface{} {
	normalized := make(map[string]interface{})
	for key, value := range atributos {
		normalized[key] = map[string]interface{}{
			"Type":  value.DataType,
			"Value": value.StringValue,
		}
	}
	return normalized
}

func PublicarNotificacion(body models.Notificacion) (msgId interface{}, outputError map[string]interface{}) {
	if usuarios, ok := body.Atributos["UsuariosDestino"].([]interface{}); ok {
		delete(body.Atributos, "UsuariosDestino")
		auxIdDeduplicacion := body.IdDeduplicacion
		for _, usuario := range usuarios {
			if idUsuario, ok := usuario.(string); ok {
				mensajeBody := body
				mensajeBody.IdDeduplicacion = auxIdDeduplicacion + idUsuario
				mensajeBody.Atributos["UsuarioDestino"] = idUsuario
				mensajeBody.IdGrupoMensaje = idUsuario
				msgId, outputError = Publicar(mensajeBody, true)
			}
		}
	} else {
		msgId, outputError = Publicar(body, false)
	}
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
		outputError = map[string]interface{}{"funcion": "/Suscribir", "err": err.Error(), "status": "502"}
		return "", outputError
	}

	client := sns.NewFromConfig(cfg)

	for _, subscriptor := range body.Suscritos {
		strAdicional := ""
		if len(subscriptor.Atributos) > 0 {
			strAdicional = ",\"" + strings.Join(subscriptor.Atributos, "\",\"") + "\""
		}
		strAtt := "{\"Destinatario\":[\"" + subscriptor.Id + "\",\"todos\"" + strAdicional + "]"
		if len(atributos) > 0 {
			for k, v := range atributos {
				strAtt += ",\"" + k + "\":[\"" + v + "\"]"
			}
		}
		strAtt += "}"
		input := &sns.SubscribeInput{
			Endpoint:              &subscriptor.Endpoint,
			Protocol:              aws.String(subscriptor.Protocolo),
			ReturnSubscriptionArn: true,
			TopicArn:              &body.ArnTopic,
			Attributes: map[string]string{
				"FilterPolicy": strAtt,
			},
		}

		result, err := client.Subscribe(context.TODO(), input)
		if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/PublicarNotificacion", "err": err.Error(), "status": "502"}
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
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.ListTopicsInput{}

	results, err := client.ListTopics(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	for _, t := range results.Topics {
		if strings.Contains(*t.TopicArn, beego.BConfig.RunMode+"-") {
			topicArn = append(topicArn, *t.TopicArn)
		}
	}
	return
}

func CrearTopic(topic models.Topic) (arn string, outputError map[string]interface{}) {
	var tags []types.Tag
	var key0 string = "name"
	var key1 string = "environment"
	var val1 string = "prod"

	val1 = beego.BConfig.RunMode

	topic.Nombre = beego.BConfig.RunMode + "-" + topic.Nombre

	if topic.EsFifo {
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
		outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err.Error(), "status": "502"}
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
			"FifoTopic":   strconv.FormatBool(topic.EsFifo),
		},
		Tags: tags,
	}

	results, err := client.CreateTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err.Error(), "status": "502"}
		return "", outputError
	}

	return *results.TopicArn, nil
}

func VerificarSuscripcion(consulta models.ConsultaSuscripcion) (suscrito bool, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/VerificarSuscripcion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ListaTopics", "err": err.Error(), "status": "502"}
		return false, outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.ListSubscriptionsByTopicInput{
		TopicArn: &consulta.ArnTopic,
	}

	results, err := client.ListSubscriptionsByTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearTopic", "err": err.Error(), "status": "502"}
		return false, outputError
	}
	logs.Debug(results)
	for _, resultado := range results.Subscriptions {
		if *resultado.Endpoint == consulta.Endpoint {
			inputSusArn := &sns.GetSubscriptionAttributesInput{
				SubscriptionArn: *&resultado.SubscriptionArn,
			}
			attributesSus, err := client.GetSubscriptionAttributes(context.TODO(), inputSusArn)
			if err == nil {
				if PendingConfirmation, err := strconv.ParseBool(attributesSus.Attributes["PendingConfirmation"]); err == nil {
					if !PendingConfirmation {
						//Esta suscrito y confirmo la suscripcion
						return true, nil
					} else {
						//Esta suscrito y no ha confirmado la suscripcion
						return false, nil
					}
				}
			}
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
		outputError = map[string]interface{}{"funcion": "/BorrarTopic", "err": err.Error(), "status": "502"}
		return outputError
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.DeleteTopicInput{
		TopicArn: &arn,
	}

	_, err = client.DeleteTopic(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarTopic", "err": err.Error(), "status": "502"}
		return outputError
	}

	return
}
