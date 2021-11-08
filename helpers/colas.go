package helpers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/udistrital/notificacion_api/models"
)

func CrearCola(cola models.Cola) (arn string, outputError map[string]interface{}) {
	var env string = "prod"

	if beego.BConfig.RunMode == "dev" || beego.BConfig.RunMode == "test" {
		env = "test"
	}

	if cola.EsFifo {
		cola.Nombre += ".fifo"
	}
	ValoresDefault(&cola)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err, "status": "502"}
		return "", outputError
	}

	client := sqs.NewFromConfig(cfg)
	policy, _ := json.Marshal(cola.Politica)
	logs.Debug(string(policy))
	input := &sqs.CreateQueueInput{
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

	result, err := client.CreateQueue(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err, "status": "502"}
		return "", outputError
	}
	return *result.QueueUrl, nil
}

func RecibirMensajes(url string) (mensajes []types.Message, outputError map[string]interface{}) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err, "status": "502"}
		return nil, outputError
	}

	client := sqs.NewFromConfig(cfg)

	input := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            &url,
		MaxNumberOfMessages: 10,
	}

	logs.Debug(url)

	result, err := client.ReceiveMessage(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err, "status": "502"}
		return nil, outputError
	}
	return result.Messages, nil
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
	logs.Debug(cola.Politica)
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
						"aws:SourceArn": "arn:aws:sns:us-east-1:*",
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
