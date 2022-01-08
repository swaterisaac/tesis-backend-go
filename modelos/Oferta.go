package modelos

import "database/sql"

type OfertaTuristica struct {
	ID           int            `json:"id"`
	Nombre       string         `json:"nombre"`
	FechaInicio  string         `json:"fecha_inicio"`
	FechaFinal   string         `json:"fecha_final"`
	Precio       sql.NullString `json:"precio"`
	Comuna       string         `json:"comuna"`
	Region       string         `json:"region"`
	Ubicacion    sql.NullString `json:"ubicacion"`
	Proveedor    string         `json:"proveedor"`
	Telefono     string         `json:"telefono"`
	Correo       string         `json:"correo"`
	Pagina       string         `json:"pagina"`
	ImagenRegion string         `json:"imagen_region"`
	ImagenOferta sql.NullString `json:"imagen_oferta"`
}

type OfertaTuristicaScrapping struct {
	OfertaTuristica
	IdComuna    int
	IdProveedor int
}
