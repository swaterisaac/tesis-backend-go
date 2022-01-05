package modelos

/*
Respuesta de elasticsearch con el formato del índice correspondiente.

Id: string del documento de elasticsearch
IdEnDb: float64 El id de la oferta turística en postgresql
Score: float64 Puntaje de relevancia del documento
FechaFinal: string Con la fecha final de cuándo vence la oferta turística
FechaInicio: string Con la fecha de inicio de validez de la oferta turística
Precio: float64 Con el precio que cuesta la oferta turística
Consideraciones: []string Lista de consideraciones que la oferta turística tiene en cuenta
Descriptor: string Con todo lo que describe a la oferta turística (Región, comuna, nombre, proveedor, etc)
*/
type RespuestaElastic struct {
	Id              string
	IdEnDb          float64
	Score           float64
	FechaFinal      string
	FechaInicio     string
	Precio          float64
	Consideraciones []string
	Descriptor      string
}

type RespuestaStatus struct {
	Code string
	Hits int
	Took int
}
