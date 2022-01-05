package busquedas

import (
	"database/sql"
	"fmt"
	"tesis/conexion"
	"tesis/modelos"
	"tesis/query"

	elasticsearch "github.com/elastic/go-elasticsearch/v6"
)

/*
Código 200: Se encontró y se retornó la oferta
Código 404: No hay oferta con esa ID
Código 500: No se pudo recuperar la oferta
*/
func ObtenerOfertaPorId(db *sql.DB, id int) (modelos.OfertaTuristica, error, int) {
	queryOferta := fmt.Sprintf(
		"SELECT ot.id, ot.nombre, ot.fecha_inicio, ot.fecha_final, ot.precio, ot.url_imagen, r.nombre, r.url_imagen, c.nombre, p.nombre, p.web, p.correo, p.telefono "+
			"FROM ofertas_turisticas ot, comunas c, proveedores p, regiones r "+
			"WHERE ot.id = %d AND ot.id_comuna = c.id AND c.id_region = r.id AND p.id = ot.id_proveedor", id)
	var oferta modelos.OfertaTuristica
	err := db.QueryRow(queryOferta).Scan(&oferta.ID, &oferta.Nombre,
		&oferta.FechaInicio, &oferta.FechaFinal, &oferta.Precio, &oferta.ImagenOferta, &oferta.Region, &oferta.ImagenRegion, &oferta.Comuna, &oferta.Proveedor,
		&oferta.Pagina, &oferta.Correo, &oferta.Telefono)
	if err == sql.ErrNoRows {
		return oferta, err, 404
	}
	if err != nil {
		return oferta, err, 500
	}
	return oferta, nil, 200
}

//Obtención de querys a elasticsearch obteniendo el comportamiento del usuario
func queryRecomendacion(db *sql.DB, idUsuario int, tamanio int, pagina int) (map[string]interface{}, error) {
	comportamiento, err := ObtenerComportamiento(db, idUsuario)
	if err != nil {
		return nil, err
	}
	queryElasticsearch := query.CrearQueryRecomendacion(comportamiento, tamanio, pagina)
	return queryElasticsearch, nil
}

func queryNormal(db *sql.DB, idUsuario int, queryUsuario string, tamanio int, pagina int) (map[string]interface{}, error) {
	comportamiento, err := ObtenerComportamiento(db, idUsuario)
	if err != nil {
		return nil, err
	}
	queryElasticsearch := query.CrearQuery(comportamiento, queryUsuario, tamanio, pagina)
	return queryElasticsearch, nil
}

//Obtención de ofertas
func obtenerOfertasGeneral(db *sql.DB, es *elasticsearch.Client, idUsuario int, queryElasticsearch map[string]interface{}) ([]modelos.OfertaTuristica, error) {

	res, r, err := conexion.BuscarRecomendacion(es, queryElasticsearch)
	if err != nil {
		return nil, err
	}
	conexion.ObtenerRespuestaStatus(res, r)
	respuestas := conexion.CastearRespuesta(r)

	var ofertas []modelos.OfertaTuristica

	for _, respuesta := range respuestas {
		oferta, err, _ := ObtenerOfertaPorId(db, int(respuesta.IdEnDb))
		if err != nil {
			return ofertas, err
		}
		ofertas = append(ofertas, oferta)
	}
	return ofertas, nil
}

func ObtenerOfertasRecomendacion(db *sql.DB, es *elasticsearch.Client, idUsuario int, tamanio int, pagina int) ([]modelos.OfertaTuristica, error) {
	queryElasticsearch, err := queryRecomendacion(db, idUsuario, tamanio, pagina)
	if err != nil {
		return nil, err
	}
	ofertas, err := obtenerOfertasGeneral(db, es, idUsuario, queryElasticsearch)
	return ofertas, err
}

func ObtenerOfertasQuery(db *sql.DB, es *elasticsearch.Client, idUsuario int, queryUsuario string, tamanio int, pagina int) ([]modelos.OfertaTuristica, error) {
	queryElasticsearch, err := queryNormal(db, idUsuario, queryUsuario, tamanio, pagina)
	if err != nil {
		return nil, err
	}
	ofertas, err := obtenerOfertasGeneral(db, es, idUsuario, queryElasticsearch)
	return ofertas, err
}

func ObtenerProveedores(db *sql.DB) ([]modelos.Region, error) {
	query := fmt.Sprintf("SELECT id, nombre FROM proveedores")
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	var proveedores []modelos.Region
	for rows.Next() {
		var proveedor modelos.Region
		if err := rows.Scan(&proveedor.ID, &proveedor.Nombre); err != nil {
			return proveedores, err
		}
		proveedores = append(proveedores, proveedor)
	}
	return proveedores, nil
}
