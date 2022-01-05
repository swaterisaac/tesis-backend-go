package modelos

import (
	"fmt"
	"strings"
)

/*
Valor: string que representa qué dato tiene el historial en concreto
Frecuencia: int que dice cuántas veces se ha repetido ese valor en el tiempo
*/
type Historial struct {
	Valor      string
	Frecuencia int
}

/*
Representa el comportamiento que es una recopilación de todos estos datos de un usuario específico
ConsideracionesMedicas: []string Lista de consideraciones médicas de un usuario
Comuna: string comuna donde vive
Region: string región donde vive
ComunaHistorial: []Historial Historial de las comunas de ofertas turísticas que ha buscado
RegiónHistorial: []Historial Historial de las regiones de ofertas turísticas que ha buscado
HistorialOfertaTuristica: []Historial Historial de las ofertas turísticas que ha buscado
HistorialBusqueda: []Historial Historial de las query o consultas que ha hecho un usuario
*/
type Comportamiento struct {
	ConsideracionesMedicas   []string
	Comuna                   string
	Region                   string
	ComunaHistorial          []Historial
	RegionHistorial          []Historial
	HistorialOfertaTuristica []Historial
	HistorialBusqueda        []Historial
}

/*
Entradas:
comportamiento: Comportamiento de un usuario

Salidas:
string: El historial de la ubicación como string

Descripción:
Convierte la lista del historial de ubicación en un string largo con todos los valores del mismo.
*/
func CastearUbicacionHistorial(comportamiento Comportamiento) string {
	comunaHistorial := ""
	for _, comuna := range comportamiento.ComunaHistorial {
		comunaHistorial += comuna.Valor + " "
	}

	regionHistorial := ""
	for _, region := range comportamiento.RegionHistorial {
		regionHistorial += region.Valor + " "
	}
	return comunaHistorial + " " + regionHistorial
}

/*
Entradas:
comportamiento: Comportamiento de un usuario

Salidas:
string: Las consideraciones como string

Descripción:
Convierte la lista de consideraciones médicas en un string largo con todos los valores del mismo.
*/
func CastearConsideraciones(comportamiento Comportamiento) string {
	consideracionesMedicas := strings.Join(comportamiento.ConsideracionesMedicas, " ")
	return consideracionesMedicas
}

/*
Entradas:
comportamiento: Comportamiento de un usuario

Salidas:
string: El historial de la ubicación como string

Descripción:
Convierte la lista del historial de ubicación en un string largo con todos los valores del mismo.
*/
func CastearOfertaHistorial(comportamiento Comportamiento) string {
	ofertaHistorial := ""
	for _, oferta := range comportamiento.HistorialOfertaTuristica {
		ofertaHistorial += oferta.Valor + " "
	}
	return ofertaHistorial
}

/*
Entradas:
comportamiento: Comportamiento de un usuario

Salidas:
string: El historial de búsqueda como string

Descripción:
Convierte la lista del historial de búsqueda en un string largo con todos los valores del mismo.
*/
func CastearBusquedaHistorial(comportamiento Comportamiento) string {
	busquedaHistorial := ""
	for _, busqueda := range comportamiento.HistorialBusqueda {
		busquedaHistorial += busqueda.Valor + " "
	}
	return busquedaHistorial
}

/*
Entradas:
comportamiento: Comportamiento de un usuario

Salidas:
string: La ubicación entera como string

Descripción:
Convierte la ubicación de un usuario en un string con la Región y Comuna
*/
func CastearUbicacion(comportamiento Comportamiento) string {
	return comportamiento.Region + " " + comportamiento.Comuna
}

/*
Entradas:
comportamiento: Comportamiento de un usuario

Salidas:
Ninguna

Descripción:
Imprime todo el comportamiento con los valores transformados a solo string.
*/
func ImprimirComportamiento(comportamiento Comportamiento) {
	fmt.Println("Consideraciones medicas: ", CastearConsideraciones(comportamiento))
	fmt.Printf("Historial ubicacion: %s \n", CastearUbicacionHistorial(comportamiento))
	fmt.Printf("Historial oferta: %s\n", CastearOfertaHistorial(comportamiento))
	fmt.Printf("Historial busqueda: %s\n", CastearBusquedaHistorial(comportamiento))
	fmt.Printf("Region y comuna: %s\n", CastearUbicacion(comportamiento))
}
