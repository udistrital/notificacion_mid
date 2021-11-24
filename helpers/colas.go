package helpers

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/udistrital/notificacion_mid/models"
)

func CrearCola(cola models.Cola) (arn string, outputError map[string]interface{}) {
	var env string = "prod"
	var fifoBool string
	var input *sqs.CreateQueueInput

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if beego.BConfig.RunMode == "dev" || beego.BConfig.RunMode == "test" {
		env = "test"
	}
	ValoresDefault(&cola)
	policy, _ := json.Marshal(cola.Politica)

	if cola.EsFifo {
		cola.Nombre += ".fifo"
		input = &sqs.CreateQueueInput{
			QueueName: &cola.Nombre,
			Attributes: map[string]string{
				"DelaySeconds":                  strconv.Itoa(cola.Retraso),
				"MessageRetentionPeriod":        strconv.Itoa(cola.Retencion),
				"MaximumMessageSize":            strconv.Itoa(cola.TamañoMaximo),
				"ReceiveMessageWaitTimeSeconds": strconv.Itoa(cola.TiempoEspera),
				"VisibilityTimeout":             strconv.Itoa(cola.EsperaVisibilidad),
				"Policy":                        string(policy),
				"FifoQueue":                     fifoBool,
			},
			Tags: map[string]string{
				"Name":        cola.Nombre,
				"Environment": env,
			},
		}
	} else {
		input = &sqs.CreateQueueInput{
			QueueName: &cola.Nombre,
			Attributes: map[string]string{
				"DelaySeconds":                  strconv.Itoa(cola.Retraso),
				"MessageRetentionPeriod":        strconv.Itoa(cola.Retencion),
				"MaximumMessageSize":            strconv.Itoa(cola.TamañoMaximo),
				"ReceiveMessageWaitTimeSeconds": strconv.Itoa(cola.TiempoEspera),
				"VisibilityTimeout":             strconv.Itoa(cola.EsperaVisibilidad),
				"Policy":                        string(policy),
			},
			Tags: map[string]string{
				"Name":        cola.Nombre,
				"Environment": env,
			},
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err, "status": "502"}
		return "", outputError
	}
	client := sqs.NewFromConfig(cfg)

	result, err := client.CreateQueue(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err, "status": "502"}
		return "", outputError
	}
	s1 := strings.Split(*result.QueueUrl, "/")
	s2 := strings.Split(s1[2], ".")
	arn = "arn:aws:sqs:" + s2[1] + ":" + s1[3] + ":" + s1[4]
	return
}

func RecibirMensajes(nombre string, tiempoOculto int) (mensajes []models.Mensaje, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RecibirMensaje", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err, "status": "502"}
		return nil, outputError
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &nombre,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err, "status": "502"}
		return nil, outputError
	}

	queueURL := resultQ.QueueUrl

	input := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   int32(tiempoOculto),
	}

	result, err := client.ReceiveMessage(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err, "status": "502"}
		return nil, outputError
	}
	for _, m := range result.Messages {
		var body map[string]interface{}
		json.Unmarshal([]byte(*m.Body), &body)
		mensajes = append(mensajes, models.Mensaje{
			Attributes:    m.Attributes,
			Body:          body,
			ReceiptHandle: *m.ReceiptHandle,
		})
	}
	return
}

func BorrarMensaje(cola string, mensaje models.Mensaje) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BorrarMensaje", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/EliminarMensaje", "err": err, "status": "502"}
		return outputError
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &cola,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/EliminarMensaje", "err": err, "status": "502"}
		return outputError
	}

	queueURL := resultQ.QueueUrl

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: &mensaje.ReceiptHandle,
	}

	_, err = client.DeleteMessage(context.TODO(), dMInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/EliminarMensaje", "err": err, "status": "502"}
		return outputError
	}

	return
}

func BorrarCola(nombre string) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err, "status": "502"}
		return outputError
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &nombre,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err, "status": "502"}
		return outputError
	}

	queueURL := resultQ.QueueUrl

	dMInput := &sqs.DeleteQueueInput{
		QueueUrl: queueURL,
	}

	_, err = client.DeleteQueue(context.TODO(), dMInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err, "status": "502"}
		return outputError
	}

	return
}

func ValoresDefault(cola *models.Cola) {
	if cola.EsperaVisibilidad == 0 {
		cola.EsperaVisibilidad = 30
	}
	if cola.Retencion == 0 {
		cola.Retencion = 345600
	}
	if cola.TamañoMaximo == 0 {
		cola.TamañoMaximo = 262144
	}
	if cola.Politica == nil {
		cola.Politica = &models.Politica{
			Version: "2008-10-17",
			Id:      "PolicySNSSQS",
			Statement: []struct {
				Sid       string
				Effect    string
				Principal struct{ AWS string }
				Action    []string
				Resource  string
				Condition map[string]map[string]string
			}{{
				Sid:      "topic-subscription-snssqs",
				Effect:   "Allow",
				Action:   []string{"sqs:SendMessage"},
				Resource: "arn:aws:sqs:*",
				Condition: map[string]map[string]string{
					"ArnLike": {
						"aws:SourceArn": "arn:aws:sns:us-east-1:*:*",
					},
				},
				Principal: struct{ AWS string }{
					AWS: "*",
				},
			},
			},
		}
	}
}
