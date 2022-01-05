package busquedas

import (
	"database/sql"
	"fmt"
	"net/http"
	"tesis/modelos"
)

func ObtenerConsideracionesMedicas(db *sql.DB) ([]modelos.Region, error, int) {
	query := fmt.Sprintf("SELECT id, nombre FROM consideraciones_medicas cm")
	rows, err := db.Query(query)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil, err, http.StatusNotFound
	}
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	var consideraciones []modelos.Region
	for rows.Next() {
		var consideracion modelos.Region
		if err := rows.Scan(&consideracion.ID, &consideracion.Nombre); err != nil {
			return consideraciones, err, http.StatusInternalServerError
		}
		consideraciones = append(consideraciones, consideracion)
	}
	return consideraciones, nil, http.StatusOK
}

func ObtenerConsideracionesPorOferta(db *sql.DB, idOferta int) ([]modelos.Region, error) {
	query := fmt.Sprintf("SELECT cm.id, cm.nombre "+
		"FROM consideraciones_medicas cm, ofertas_turisticas ot, oferta_consideraciones oc "+
		"WHERE ot.id = %d AND ot.id = oc.id_oferta AND cm.id = oc.id_consideracion", idOferta)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var consideraciones []modelos.Region
	for rows.Next() {
		var consideracion modelos.Region
		if err := rows.Scan(&consideracion.ID, &consideracion.Nombre); err != nil {
			return consideraciones, err
		}
		consideraciones = append(consideraciones, consideracion)
	}
	return consideraciones, nil
}
