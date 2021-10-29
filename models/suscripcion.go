package models

type Suscripcion struct {
	//Lista de endpoints a suscribir
	Suscritos []struct {
		//Metadato usado para el filtrado de mensajes, es decir, para que un suscriptor sólo reciba
		//notificaciones destinadas a él.
		Id string

		//Puede ser un correo (email), número telefónico (sms), endpoint (application), o un ARN de un
		//servicio AWS.
		Endpoint string

		//Para el caso de topics estándar, se suscriben endpoints con protocolo
		//kinesis, lambda, sqs, email, email-json, http, https, aplication, sms o firehouse.
		//Para más información consultar https://docs.aws.amazon.com/sns/latest/api/API_Subscribe.html

		//En el caso de topics tipo FIFO, sólo se debe utilizar el valor sqs para el protocolo
		Protocolo string
	}

	//ARN de topic al que se quiere suscribir a la lista
	ArnTopic string
}
