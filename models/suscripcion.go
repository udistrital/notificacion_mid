package models

type Suscripcion struct {
	//Lista de endpoints a suscribir
	Suscritos []Endpoint

	//ARN de topic al que se quiere suscribir a la lista
	ArnTopic string
}
