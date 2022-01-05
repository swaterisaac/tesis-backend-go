package conexion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"tesis/modelos"

	elasticsearch "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

func BuscarRecomendacion(es *elasticsearch.Client, query map[string]interface{}) (*esapi.Response, map[string]interface{}, error) {
	var buf bytes.Buffer
	var r map[string]interface{}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, r, err
	}
	res, err := es.Search(
		es.Search.WithIndex("testlogstash"),
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, r, err
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, r, err
	}
	return res, r, err
}

func CastearRespuesta(r map[string]interface{}) []modelos.RespuestaElastic {
	var respuestas []modelos.RespuestaElastic
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var respuesta modelos.RespuestaElastic
		respuesta.Id = fmt.Sprintf("%v", hit.(map[string]interface{})["_id"])
		respuesta.IdEnDb = hit.(map[string]interface{})["_source"].(map[string]interface{})["idEnDb"].(float64)
		respuesta.Descriptor = hit.(map[string]interface{})["_source"].(map[string]interface{})["descriptor"].(string)
		respuesta.FechaFinal = hit.(map[string]interface{})["_source"].(map[string]interface{})["fecha_final"].(string)
		respuesta.FechaInicio = hit.(map[string]interface{})["_source"].(map[string]interface{})["fecha_inicio"].(string)
		respuesta.Precio = hit.(map[string]interface{})["_source"].(map[string]interface{})["precio"].(float64)
		respuesta.Score = hit.(map[string]interface{})["_score"].(float64)
		consideraciones := hit.(map[string]interface{})["_source"].(map[string]interface{})["consideraciones"].([]interface{})
		//Castear arreglo de []interfaz{} a []string
		respuesta.Consideraciones = make([]string, len(consideraciones))
		for i := range consideraciones {
			respuesta.Consideraciones[i] = consideraciones[i].(string)
		}
		respuestas = append(respuestas, respuesta)

		//log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"].(map[string]interface{})["descriptor"])
	}
	return respuestas
}

func ObtenerRespuestaStatus(res *esapi.Response, r map[string]interface{}) modelos.RespuestaStatus {
	var status modelos.RespuestaStatus
	status.Code = res.Status()
	status.Hits = int(r["hits"].(map[string]interface{})["total"].(float64))
	status.Took = int(r["took"].(float64))
	return status
}
