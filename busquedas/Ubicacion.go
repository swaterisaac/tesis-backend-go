package busquedas

import (
	"database/sql"
	"fmt"
	"tesis/modelos"
)

func ObtenerRegiones(db *sql.DB) ([]modelos.Region, error, int) {
	query := fmt.Sprintf("SELECT id, nombre FROM regiones")
	rows, err := db.Query(query)

	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil, nil, 404
	}
	if err != nil {
		return nil, err, 500
	}
	var regiones []modelos.Region
	for rows.Next() {
		var region modelos.Region
		if err := rows.Scan(&region.ID, &region.Nombre); err != nil {
			return regiones, err, 500
		}
		regiones = append(regiones, region)
	}
	return regiones, nil, 200
}

func ObtenerComunasPorRegion(db *sql.DB, idRegion int) ([]modelos.Region, error, int) {
	query := fmt.Sprintf("SELECT c.id, c.nombre FROM comunas c, regiones r WHERE c.id_region = %d AND r.id = c.id_region", idRegion)
	rows, err := db.Query(query)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil, nil, 404
	}
	if err != nil {
		return nil, err, 500
	}
	var comunas []modelos.Region
	for rows.Next() {
		var comuna modelos.Region
		if err := rows.Scan(&comuna.ID, &comuna.Nombre); err != nil {
			return comunas, err, 500
		}
		comunas = append(comunas, comuna)
	}
	return comunas, nil, 200
}
