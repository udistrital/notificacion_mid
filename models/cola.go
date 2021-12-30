package models

type Cola struct {
	NombreCola        string
	EsFifo            bool
	EsperaVisibilidad int
	Retencion         int
	Retraso           int
	TamañoMaximo      int
	TiempoEspera      int
	Politica          *Politica
}
