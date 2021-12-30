package models

type Notificacion struct {
	//Id del remitente para identificación y posibilidad de respuesta
	RemitenteId string

	//Lista de ids de destinatarios para la difusión de la notificación
	DestinatarioId []string

	//Asunto del mensaje para notificaciones por correo
	Asunto string

	//Mensaje de la notificación a enviar
	Mensaje string

	//Metadatos de la notificación, puede usarse para filtrado de
	//información en sns o para transferir datos entre los sistemas,
	//para más información visitar (https://docs.aws.amazon.com/es_es/sns/latest/dg/sns-message-attributes.html)
	Atributos map[string]interface{}

	//Identifica el ARN del topic creado en SNS, para más
	//información sobre los topics visitar (https://docs.aws.amazon.com/sns/latest/dg/sns-create-topic.html)
	ArnTopic string

	//En colas fifo, se identifica cada mensaje con un id para que no se repita la misma notificacion varias veces,
	//comparando los id de deduplicacion
	IdDeduplicacion string

	//En colas fifo, los mensajes se separan por grupos, y los mensajes de cada grupo son entregados en el orden
	//en que fueron enviados
	IdGrupoMensaje string
}
