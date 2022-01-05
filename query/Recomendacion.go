package query

import "tesis/modelos"

//Puntajes de cada campo del comportamiento
//Puntaje consideraciones medicas
const puntajeConsideraciones float32 = 0.35

//Puntaje de las ubicaciones del historial de las ofertas buscadas
const puntajeUbicacionHistorial float32 = 0.02

//Puntaje del historial de las ofertas turísticas
const puntajeOfertasHistorial float32 = 0.08

//Puntaje del historial de consultas hechas por el usuario
const puntajeBusquedasHistorial float32 = 0.02

//Puntaje NEGATIVO de la ubicacion del usuario
//const puntajeUbicacion float32 = 0.1

//Para modificar cuántas comunas, regiones y ofertas se contarán en el historial, modificar
//el archivo de busquedas/comportamiento

//Para modificar el límite de historial búsqueda a guardar, está en el archivo busquedas/Historial.go

/*
Entradas:
comportamiento: modelos.Comportamiento El comportamiento de un usuario
puntaje: float32 el puntaje que se le asigna a esta micro-consulta

Salidas:
[]map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Retorna la micro-consulta de las consideraciones médicas para luego sumarla al JSON general.
*/
func rescatarConsideraciones(comportamiento modelos.Comportamiento, puntaje float32) []map[string]interface{} {
	var consideraciones []map[string]interface{}
	for _, consideracion := range comportamiento.ConsideracionesMedicas {
		interfaz := map[string]interface{}{
			"match": map[string]interface{}{
				"consideraciones": map[string]interface{}{
					"query": consideracion,
					"boost": puntaje,
				},
			},
		}
		consideraciones = append(consideraciones, interfaz)
	}

	return consideraciones
}

/*
Entradas:
comportamiento: modelos.Comportamiento El comportamiento de un usuario
puntaje: float32 el puntaje que se le asigna a esta micro-consulta

Salidas:
[]map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Retorna la micro-consulta del historial de ubicación para luego sumarla al JSON general.
*/
func rescatarUbicacionHistorial(comportamiento modelos.Comportamiento, puntaje float32) map[string]interface{} {
	ubicacionHistorial := map[string]interface{}{
		"match": map[string]interface{}{
			"descriptor": map[string]interface{}{
				"query": modelos.CastearUbicacionHistorial(comportamiento),
				"boost": puntaje,
			},
		},
	}
	return ubicacionHistorial
}

/*
Entradas:
comportamiento: modelos.Comportamiento El comportamiento de un usuario
puntaje: float32 el puntaje que se le asigna a esta micro-consulta

Salidas:
[]map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Retorna la micro-consulta del historial de ofertas turísticas para luego sumarla al JSON general.
*/
func rescatarOfertasHistorial(comportamiento modelos.Comportamiento, puntaje float32) map[string]interface{} {
	ofertasHistorial := map[string]interface{}{
		"match": map[string]interface{}{
			"descriptor": map[string]interface{}{
				"query": modelos.CastearOfertaHistorial(comportamiento),
				"boost": puntaje,
			},
		},
	}
	return ofertasHistorial
}

/*
Entradas:
comportamiento: modelos.Comportamiento El comportamiento de un usuario
puntaje: float32 el puntaje que se le asigna a esta micro-consulta

Salidas:
[]map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Retorna la micro-consulta del historial de queries para luego sumarla al JSON general.
*/
func rescatarBusquedaHistorial(comportamiento modelos.Comportamiento, puntaje float32) map[string]interface{} {
	busquedaHistorial := map[string]interface{}{
		"match": map[string]interface{}{
			"descriptor": map[string]interface{}{
				"query": modelos.CastearBusquedaHistorial(comportamiento),
				"boost": puntaje,
			},
		},
	}
	return busquedaHistorial
}

/*
Entradas:
comportamiento: modelos.Comportamiento El comportamiento de un usuario
puntaje: float32 el puntaje que se le asigna a esta micro-consulta

Salidas:
[]map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Retorna la micro-consulta de la ubicación del usuario para luego sumarla al JSON general.
*/
func rescatarUbicacion(comportamiento modelos.Comportamiento, puntaje float32) map[string]interface{} {
	ubicacion := map[string]interface{}{
		"match": map[string]interface{}{
			"descriptor": map[string]interface{}{
				"query": modelos.CastearUbicacion(comportamiento),
				"boost": puntaje,
			},
		},
	}
	return ubicacion
}

/*
Entradas:
comportamiento: modelos.Comportamiento El comportamiento de un usuario

Salidas:
[]map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Suma toda las micro-consultas para crear el JSON de búsqueda a elasticsearch
*/
func CrearQueryRecomendacion(comportamiento modelos.Comportamiento, tamanio int, pagina int) map[string]interface{} {
	desde := tamanio * (pagina - 1)
	clausulaShould := []interface{}{}
	consideraciones := rescatarConsideraciones(comportamiento, puntajeConsideraciones)
	for _, consideracion := range consideraciones {
		clausulaShould = append(clausulaShould, consideracion)
	}
	clausulaShould = append(clausulaShould, rescatarUbicacionHistorial(comportamiento, puntajeUbicacionHistorial))
	clausulaShould = append(clausulaShould, rescatarOfertasHistorial(comportamiento, puntajeOfertasHistorial))
	clausulaShould = append(clausulaShould, rescatarBusquedaHistorial(comportamiento, puntajeBusquedasHistorial))
	query := map[string]interface{}{
		"size": tamanio,
		"from": desde,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": clausulaShould,
			},
		},
	}
	return query
}
