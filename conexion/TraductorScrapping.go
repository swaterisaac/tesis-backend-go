package conexion

import (
	"database/sql"
	"tesis/modelos"
)

func TraductorScrappingGeneral(db *sql.DB, rows *sql.Rows) ([]modelos.OfertaTuristica, error) {

	var ofertas []modelos.OfertaTuristica

	for rows.Next() {
		var oferta modelos.OfertaTuristica
		if err := rows.Scan(&oferta.ID, &oferta.Comuna, &oferta.FechaFinal,
			&oferta.FechaInicio, &oferta.Nombre, &oferta.Precio,
			&oferta.Proveedor, &oferta.Region); err != nil {
			return ofertas, err
		}
		ofertas = append(ofertas, oferta)
	}

	if err := rows.Err(); err != nil {
		return ofertas, err
	}
	return ofertas, nil
}
