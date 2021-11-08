package models

type Cola struct {
	Nombre            string
	EsFifo            bool
	EsperaVisibilidad int
	Retencion         int
	Retraso           int
	TamañoMaximo      int
	TiempoEspera      int
	Politica          *Politica
	ArnSns            string
}
