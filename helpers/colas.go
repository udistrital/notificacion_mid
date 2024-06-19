package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/udistrital/notificacion_mid/models"
)

func CrearCola(cola models.Cola) (arn string, outputError map[string]interface{}) {
	var env string = beego.BConfig.RunMode
	fifoBool := strconv.FormatBool(cola.EsFifo)
	var input *sqs.CreateQueueInput

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	ValoresDefault(&cola)
	policy, _ := json.Marshal(cola.Politica)

	if cola.EsFifo {
		cola.NombreCola += ".fifo"
		input = &sqs.CreateQueueInput{
			QueueName: &cola.NombreCola,
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
				"Name":        cola.NombreCola,
				"Environment": env,
			},
		}
	} else {
		input = &sqs.CreateQueueInput{
			QueueName: &cola.NombreCola,
			Attributes: map[string]string{
				"DelaySeconds":                  strconv.Itoa(cola.Retraso),
				"MessageRetentionPeriod":        strconv.Itoa(cola.Retencion),
				"MaximumMessageSize":            strconv.Itoa(cola.TamañoMaximo),
				"ReceiveMessageWaitTimeSeconds": strconv.Itoa(cola.TiempoEspera),
				"VisibilityTimeout":             strconv.Itoa(cola.EsperaVisibilidad),
				"Policy":                        string(policy),
			},
			Tags: map[string]string{
				"Name":        cola.NombreCola,
				"Environment": env,
			},
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err.Error(), "status": "502"}
		return "", outputError
	}
	client := sqs.NewFromConfig(cfg)

	result, err := client.CreateQueue(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/CrearCola", "err": err.Error(), "status": "502"}
		return "", outputError
	}
	s1 := strings.Split(*result.QueueUrl, "/")
	s2 := strings.Split(s1[2], ".")
	arn = "arn:aws:sqs:" + s2[1] + ":" + s1[3] + ":" + s1[4]
	return
}

func RecibirMensajes(nombre string, tiempoOculto int, numMax int) (mensajes []models.Mensaje, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RecibirMensaje", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	nombre = beego.BConfig.RunMode + "-" + nombre

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &nombre,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	queueURL := resultQ.QueueUrl

	input := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: int32(numMax),
		VisibilityTimeout:   int32(tiempoOculto),
		WaitTimeSeconds:     *aws.Int32(5),
	}

	result, err := client.ReceiveMessage(context.TODO(), input)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajes", "err": err.Error(), "status": "502"}
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

func RecibirMensajesPorUsuario(nombre string, id_usuario string, numRevisados int, idMensaje string) (mensajes []models.Mensaje, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/RecibirMensajesPorUsuario", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	nombreAWS := beego.BConfig.RunMode + "-" + nombre

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajesPorUsuario", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &nombreAWS,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/RecibirMensajesPorUsuario", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	queueURL := resultQ.QueueUrl

	var listaTodosMensajes []models.Mensaje
	var listaPendientes []models.Mensaje
	var listaRevisados []models.Mensaje

	input := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   3,
		WaitTimeSeconds:     *aws.Int32(5),
	}

	var auxMenRef models.Mensaje
	for {
		// Obtener mensajes
		result, err := client.ReceiveMessage(context.TODO(), input)
		if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RecibirMensajesPorUsuario", "err": err.Error(), "status": "502"}
			return nil, outputError
		}

		if len(result.Messages) == 0 {
			break
		}

		var loteMensajes []models.Mensaje

		// Crear lote de mensajes  y guardado
		for _, m := range result.Messages {
			var body map[string]interface{}
			json.Unmarshal([]byte(*m.Body), &body)
			mensaje := models.Mensaje{Attributes: m.Attributes, Body: body, ReceiptHandle: *m.ReceiptHandle}
			loteMensajes = append(loteMensajes, mensaje)

			atributosMensaje := mensaje.Body["MessageAttributes"].(map[string]interface{})
			usuarioDestino, usuarioDestinoOk := atributosMensaje["UsuarioDestino"].(map[string]interface{})
			estadoMensaje, estadoMensajeOk := atributosMensaje["EstadoMensaje"].(map[string]interface{})

			if !usuarioDestinoOk || !estadoMensajeOk {
				outputError = map[string]interface{}{"funcion": "/RecibirMensajesPorUsuario", "err": "Error en estructura de datos", "status": "400"}
				return nil, outputError
			}

			// Cambiar el estado de un mensaje a revisado buscando por IdReferencia, o de lo contrario por MessageId
			refId, refIdOk := atributosMensaje["IdReferencia"].(map[string]interface{})
			msgId, msgIdOk := mensaje.Body["MessageId"].(string)
			if idMensaje != "" {
				if refIdOk && idMensaje == refId["Value"].(string) || msgIdOk && idMensaje == msgId {
					estadoMensaje["Value"] = "revisado"
					auxMenRef = mensaje
				} else {
					listaTodosMensajes = append(listaTodosMensajes, mensaje)
				}
			} else {
				listaTodosMensajes = append(listaTodosMensajes, mensaje)
			}

			// Guardar mensajes(pendientes y revisados), si corresponde al identificador de un usuario
			if usuarioDestino["Value"] == id_usuario {
				if estadoMensaje["Value"] == "pendiente" {
					listaPendientes = append([]models.Mensaje{mensaje}, listaPendientes...)
				} else if estadoMensaje["Value"] == "revisado" {
					if refIdOk && refId["Value"] != idMensaje {
						listaRevisados = append([]models.Mensaje{mensaje}, listaRevisados...)
					}
				}
			}
		}

		// Eliminar lote de mensajes
		entries := make([]types.DeleteMessageBatchRequestEntry, len(loteMensajes))
		for msgIndex := range loteMensajes {
			entries[msgIndex].Id = aws.String(fmt.Sprintf("%v", msgIndex))
			entries[msgIndex].ReceiptHandle = &loteMensajes[msgIndex].ReceiptHandle
		}

		_, err = client.DeleteMessageBatch(context.TODO(), &sqs.DeleteMessageBatchInput{
			Entries:  entries,
			QueueUrl: queueURL,
		})
		if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/RecibirMensajesPorUsuario", "err": err.Error(), "status": "502"}
			return nil, outputError
		}
	}

	// Colocar el mensaje cambiado a revisado al inicio de la lista de revisados
	// y al final de todos los mensajes, para que al registrar y obtener nuevamente este en la primera posición
	if idMensaje != "" && len(auxMenRef.Body) != 0 {
		listaRevisados = append([]models.Mensaje{auxMenRef}, listaRevisados...)
		listaTodosMensajes = append(listaTodosMensajes, auxMenRef)
	}

	// Filtrar por mensajes en estado pendiente y luego por mensajes en estado revisado
	mensajes = append(mensajes, listaPendientes...)
	if numRevisados > len(listaRevisados) || numRevisados == -1 {
		mensajes = append(mensajes, listaRevisados...)
	} else {
		mensajes = append(mensajes, listaRevisados[:numRevisados]...)
	}

	//Registrar nuevamente todos los mensajes
	tamBloque := 10
	for i := 0; i < len(listaTodosMensajes); i += tamBloque {
		fin := i + tamBloque
		if fin > len(listaTodosMensajes) {
			fin = len(listaTodosMensajes)
		}
		subLista := listaTodosMensajes[i:fin]

		// Publicar notificaciones por lote de 10 mensajes
		if errPublish := PublicarLote(subLista); errPublish != nil {
			logs.Error(errPublish)
		}
	}
	return mensajes, nil
}

func EsperarMensajes(nombre string, tiempoEspera int, cantidad int, filtro map[string]string) (mensajes []models.Mensaje, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/EsperarMensajes", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	nombre = beego.BConfig.RunMode + "-" + nombre

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/EsperarMensajes", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	client := sqs.NewFromConfig(cfg)

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &nombre,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/EsperarMensajes", "err": err.Error(), "status": "502"}
		return nil, outputError
	}

	queueURL := resultQ.QueueUrl

	restante := tiempoEspera
	var mensajesRecibidos []models.Mensaje
	for restante > 0 {
		input := &sqs.ReceiveMessageInput{
			MessageAttributeNames: []string{
				string(types.QueueAttributeNameAll),
			},
			QueueUrl:            queueURL,
			MaxNumberOfMessages: 10,
			VisibilityTimeout:   int32(restante),
			WaitTimeSeconds:     0,
		}
		result, err := client.ReceiveMessage(context.TODO(), input)
		if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/EsperarMensajes", "err": err.Error(), "status": "502"}
			return nil, outputError
		}
		for _, m := range result.Messages {
			var body map[string]interface{}
			json.Unmarshal([]byte(*m.Body), &body)
			mensajesRecibidos = append(mensajesRecibidos, models.Mensaje{
				Attributes:    m.Attributes,
				Body:          body,
				ReceiptHandle: *m.ReceiptHandle,
			})
		}
		time.Sleep(time.Second)
		restante--
	}
	for _, m := range mensajesRecibidos {
		atributos := m.Body["MessageAttributes"].(map[string]interface{})
		and := false
		for k := range atributos {
			if (atributos[k].(map[string]interface{})["Type"].(string) != "String" && strings.Contains(atributos[k].(map[string]interface{})["Value"].(string), "\""+filtro[k]+"\"")) ||
				(atributos[k].(map[string]interface{})["Type"].(string) == "String" && atributos[k].(map[string]interface{})["Value"].(string) == filtro[k]) ||
				filtro[k] == "" {
				and = true
			} else {
				and = false
				break
			}
		}
		if and {
			mensajes = append(mensajes, m)
		}
	}

	if cantidad != 0 && cantidad < len(mensajes) {
		mensajes = mensajes[:cantidad]
	}
	time.Sleep(time.Second)
	return
}

func BorrarMensaje(cola string, mensaje models.Mensaje) (conteo int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BorrarMensaje", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarMensaje", "err": err.Error(), "status": "502"}
		return 0, outputError
	}

	client := sqs.NewFromConfig(cfg)

	cola = beego.BConfig.RunMode + "-" + cola

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &cola,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarMensaje", "err": err.Error(), "status": "502"}
		return 0, outputError
	}

	queueURL := resultQ.QueueUrl

	err1 := Eliminar(queueURL, mensaje, client)
	if err1 != nil {
		logs.Error(err1)
		outputError = map[string]interface{}{"funcion": "/BorrarMensaje", "err": err1, "status": "502"}
		return 0, outputError
	}
	conteo++

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
		outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err.Error(), "status": "502"}
		return outputError
	}

	client := sqs.NewFromConfig(cfg)

	nombre = beego.BConfig.RunMode + "-" + nombre

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &nombre,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err.Error(), "status": "502"}
		return outputError
	}

	queueURL := resultQ.QueueUrl

	dMInput := &sqs.DeleteQueueInput{
		QueueUrl: queueURL,
	}

	_, err = client.DeleteQueue(context.TODO(), dMInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarCola", "err": err.Error(), "status": "502"}
		return outputError
	}

	return
}

func ValoresDefault(cola *models.Cola) {
	cola.NombreCola = beego.BConfig.RunMode + "-" + cola.NombreCola
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
						"aws:SourceArn": "arn:aws:sns:us-east-1:*:" + beego.BConfig.RunMode + "-*",
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

func BorrarMensajeFiltro(filtro models.Filtro) (conteo int, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/BorrarMensajeFiltro", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var mensajes []models.Mensaje

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarMensajeFiltro", "err": err.Error(), "status": "502"}
		return 0, outputError
	}

	client := sqs.NewFromConfig(cfg)

	filtro.NombreCola = beego.BConfig.RunMode + "-" + filtro.NombreCola

	qUInput := &sqs.GetQueueUrlInput{
		QueueName: &filtro.NombreCola,
	}

	resultQ, err := client.GetQueueUrl(context.TODO(), qUInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/BorrarMensajeFiltro", "err": err.Error(), "status": "502"}
		return 0, outputError
	}

	queueURL := resultQ.QueueUrl

	input := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   2,
	}
	for {
		result, err := client.ReceiveMessage(context.TODO(), input, func(o *sqs.Options) {
			o.ClientLogMode = aws.LogRequestWithBody | aws.LogResponseWithBody
		})
		if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/BorrarMensajeFiltro", "err": err.Error(), "status": "502"}
			return 0, outputError
		}
		if len(result.Messages) > 0 {
			for _, m := range result.Messages {
				var body map[string]interface{}
				json.Unmarshal([]byte(*m.Body), &body)
				mensajes = append(mensajes, models.Mensaje{
					Attributes:    m.Attributes,
					Body:          body,
					ReceiptHandle: *m.ReceiptHandle,
				})
			}
		} else {
			break
		}
	}

	for _, m := range mensajes {
		atributos := m.Body["MessageAttributes"].(map[string]interface{})
		for key, value := range atributos {
			if ContainsJson(filtro.Filtro[key], value.(map[string]interface{})) || ContainsString(filtro.Filtro[key], "All") {
				err1 := Eliminar(queueURL, m, client)
				if err1 != nil {
					logs.Error(err1)
					outputError = map[string]interface{}{"funcion": "/BorrarMensajeFiltro", "err": err1, "status": "502"}
					return 0, outputError
				}
				conteo++
			}
		}
	}
	time.Sleep(2 * time.Second)
	return
}

func Eliminar(urlCola *string, mensaje models.Mensaje, client *sqs.Client) (outputError map[string]interface{}) {
	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      urlCola,
		ReceiptHandle: &mensaje.ReceiptHandle,
	}

	_, err := client.DeleteMessage(context.TODO(), dMInput)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/EliminarMensaje", "err": err.Error(), "status": "502"}
		return outputError
	}
	return
}
