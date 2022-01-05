package query

import "tesis/modelos"

//Puntaje que se asigna a la query de un usuario
const puntajeQuery float32 = 0.8

/*
Entradas:
queryEntrada: string con la query de búsqueda del usuario

Salidas:
map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Hace la parte de la query para plasmarla en una consulta JSON a elasticsearch.
*/
func rescatarQuery(queryEntrada string) map[string]interface{} {
	query := map[string]interface{}{
		"match": map[string]interface{}{
			"descriptor": map[string]interface{}{
				"query": queryEntrada,
				"boost": puntajeQuery,
			},
		},
	}
	return query
}

/*
Entradas:
comportamiento: modelos.Comportamiento Con el comportamiento del usuario
queryEntrada: string Con la query o consulta que ha puesto el usuario

Salidas:
map[string]interface{}: Casteo de JSON para la consulta en elasticsearch

Descripción:
Recopila todas las consideraciones más la query correspondiente y la devuelve como JSON para
elasticsearch
*/
func CrearQuery(comportamiento modelos.Comportamiento, queryEntrada string, tamanio int, pagina int) map[string]interface{} {
	desde := tamanio * (pagina - 1)

	const reduccion float32 = 0.4
	puntajeBusquedasHistorialQ := puntajeBusquedasHistorial * reduccion
	puntajeConsideracionesQ := puntajeConsideraciones * reduccion
	puntajeOfertasHistorialQ := puntajeOfertasHistorial * reduccion
	puntajeUbicacionHistorialQ := puntajeUbicacionHistorial * reduccion

	clausulaShould := []interface{}{}
	consideraciones := rescatarConsideraciones(comportamiento, puntajeConsideracionesQ)
	for _, consideracion := range consideraciones {
		clausulaShould = append(clausulaShould, consideracion)
	}
	clausulaShould = append(clausulaShould, rescatarUbicacionHistorial(comportamiento, puntajeUbicacionHistorialQ))
	clausulaShould = append(clausulaShould, rescatarOfertasHistorial(comportamiento, puntajeOfertasHistorialQ))
	clausulaShould = append(clausulaShould, rescatarBusquedaHistorial(comportamiento, puntajeBusquedasHistorialQ))
	clausulaShould = append(clausulaShould, rescatarQuery(queryEntrada))
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
