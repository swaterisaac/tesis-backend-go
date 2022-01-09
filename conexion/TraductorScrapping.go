package conexion

import (
	"database/sql"
	"tesis/modelos"
)

func TraductorScrappingGeneral(db *sql.DB, rows *sql.Rows) ([]modelos.OfertaTuristicaScrapping, error) {
	var auxScan sql.NullString
	var ofertas []modelos.OfertaTuristicaScrapping

	for rows.Next() {
		var oferta modelos.OfertaTuristicaScrapping
		if err := rows.Scan(&oferta.ID, &oferta.Comuna, &oferta.FechaFinal,
			&oferta.FechaInicio, &oferta.Nombre, &oferta.Precio,
			&oferta.Proveedor, &oferta.Region); err != nil {
			return ofertas, err
		}
		ofertas = append(ofertas, oferta)
	}

	if err := rows.Scan(&auxScan); err != nil {
		return ofertas, err
	}

	return ofertas, nil
}

func TraductorScrappingChilepass(db *sql.DB, rows *sql.Rows) ([]modelos.OfertaTuristicaScrapping, error) {
	var auxScan sql.NullString
	var ofertas []modelos.OfertaTuristicaScrapping

	for rows.Next() {
		var oferta modelos.OfertaTuristicaScrapping
		if err := rows.Scan(&oferta.ID, &oferta.Nombre, &oferta.IdComuna, &oferta.FechaInicio,
			&oferta.Precio, &oferta.IdProveedor, &oferta.Ubicacion); err != nil {
			return ofertas, err
		}
		ofertas = append(ofertas, oferta)
	}

	if err := rows.Scan(&auxScan); err != nil {
		return ofertas, err
	}
	return ofertas, nil
}
