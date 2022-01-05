package busquedas

import (
	"database/sql"
	"sync"
	"tesis/modelos"
)

//Cuántas regiones se contarán en la recomendación
const topRegiones = 10

//Cuántas comunas se contarán en la recomendación
const topComunas = 10

//Cuántas ofertas se contarán en la recomendación
const topOfertas = 5

/*
Entradas:
db: *sql.DB sacado de la librería database/sql. Manejador de conexión a bases de datos mysql
(pgsql en este caso)
idUsuario: int que indica el id del usuario a sacar los datos.

Salidas:
modelos.Comportamiento: sacado de un struct en el paquete modelos que representa todo el comportamiento
recopilado de un usuario en la base de datos.
error: Error en caso de que algo salga mal

Descripción:
Función que usa todas las funciones con goroutines para recopilar el comportamiento de un usuario y
lo retorna en forma de struct.
*/
func ObtenerComportamiento(db *sql.DB, idUsuario int) (modelos.Comportamiento, error) {
	canalError := make(chan error)

	canalContenidoHistorialBusqueda := make(chan []modelos.Historial)
	canalContenidoHistorialOferta := make(chan []modelos.Historial)
	canalContenidoHistorialRegion := make(chan []modelos.Historial)
	canalContenidoHistorialComuna := make(chan []modelos.Historial)
	canalContenidoRegion := make(chan string)
	canalContenidoComuna := make(chan string)
	canalContenidoConsideraciones := make(chan []string)

	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		wg.Wait()
		close(canalError)
	}()

	go func() {
		historialBusqueda, err := ObtenerHistorialBusquedaUsuario(db, idUsuario)
		if err != nil {
			canalError <- err

		}
		wg.Done()
		canalContenidoHistorialBusqueda <- historialBusqueda
	}()
	go func() {
		historialOfertas, err := ObtenerHistorialOfertasUsuario(db, idUsuario, topOfertas)
		if err != nil {
			canalError <- err
		}
		wg.Done()
		canalContenidoHistorialOferta <- historialOfertas
	}()
	go func() {
		historialRegion, err := ObtenerHistorialRegionUsuario(db, idUsuario, topRegiones)
		if err != nil {
			canalError <- err
		}
		wg.Done()
		canalContenidoHistorialRegion <- historialRegion
	}()
	go func() {
		historialComuna, err := ObtenerHistorialComunaUsuario(db, idUsuario, topComunas)
		if err != nil {
			canalError <- err
		}
		wg.Done()
		canalContenidoHistorialComuna <- historialComuna
	}()
	go func() {
		comuna, region, err := ObtenerUbicacionUsuario(db, idUsuario)
		if err != nil {
			canalError <- err
		}
		wg.Done()
		canalContenidoComuna <- comuna
		canalContenidoRegion <- region
	}()
	go func() {

		consideraciones, err := ObtenerConsideracionesMedicasUsuario(db, idUsuario)
		if err != nil {
			canalError <- err
		}
		wg.Done()
		canalContenidoConsideraciones <- consideraciones
	}()
	var comportamiento modelos.Comportamiento
	//Manejo de errores
	for err := range canalError {
		if err != nil {
			return comportamiento, err
		}
	}

	comportamiento.Comuna = <-canalContenidoComuna
	comportamiento.Region = <-canalContenidoRegion
	comportamiento.HistorialBusqueda = <-canalContenidoHistorialBusqueda
	comportamiento.HistorialOfertaTuristica = <-canalContenidoHistorialOferta
	comportamiento.ComunaHistorial = <-canalContenidoHistorialComuna
	comportamiento.RegionHistorial = <-canalContenidoHistorialRegion
	comportamiento.ConsideracionesMedicas = <-canalContenidoConsideraciones

	return comportamiento, nil

}
