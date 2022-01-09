package conexion

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tesis/modelos"
	"time"
)

type traductor func(db *sql.DB, rows *sql.Rows) ([]modelos.OfertaTuristicaScrapping, error)
type queryNormal func(int) string
type queryCreacion func(modelos.OfertaTuristicaScrapping) string

func queryActualizarScrapping(ultimoId int, nombreDBScrapping string) string {
	return fmt.Sprintf("UPDATE scrapings SET ultimatupla = %d WHERE nombredb = '%s' returning id", ultimoId, nombreDBScrapping)
}

func ListenerGeneral(dbApp *sql.DB, dbScrapping *sql.DB, funcionQuery queryNormal, funcionTraductor traductor, funcionCreacion queryCreacion, nombreDBScrapping string) {
	const segundosDetenerError = 120
	const segundosDetenerFinal = 240
	var auxScan sql.NullString
	for {
		//¿Existe la tupla de ese scraping?
		var existeTabla bool
		queryExisteTabla := fmt.Sprintf("select exists (select 1 from scrapings s where s.nombredb = '%s')", nombreDBScrapping)
		err := dbApp.QueryRow(queryExisteTabla).Scan(&existeTabla)
		if err != nil {
			//TODO: Repetir en loop este error
			log.Println("Hay un error al capturar la veracidad de una tupla\nError: ", err)
			time.Sleep(segundosDetenerError * time.Second)
			continue
		}
		//Caso de que no exista la tupla, se crea
		if !existeTabla {
			log.Printf("No existe la tupla %s, creando...\n", nombreDBScrapping)
			queryCrearTupla := fmt.Sprintf("INSERT INTO scrapings(nombredb) values ('%s')", nombreDBScrapping)
			err := dbApp.QueryRow(queryCrearTupla).Err()
			if err != nil {
				//TODO: Repetir en loop este error
				log.Println("Hay un error al crear la nueva tupla\nError: ", err)
				time.Sleep(segundosDetenerError * time.Second)
				continue
			}
			log.Printf("Se ha creado la tupla %s", nombreDBScrapping)
		}

		//Obtener el valor de la última tupla para sacar datos desde esa ID
		var ultimaTupla int
		queryUltimaTupla := fmt.Sprintf("SELECT s.ultimatupla FROM scrapings s WHERE s.nombredb = '%s'", nombreDBScrapping)
		err = dbApp.QueryRow(queryUltimaTupla).Scan(&ultimaTupla)
		if err != nil {
			//TODO: Repetir en loop este error
			log.Println("Error al obtener el valor de última tupla\nError: ", err)
			time.Sleep(segundosDetenerError * time.Second)
			continue
		}
		log.Printf("Se obtuvo el valor de la última tupla: %d\n", ultimaTupla)

		//Se obtienen las ofertas de la base de datos de scrapping
		log.Println("Obteniendo las últimas ofertas...")
		query := funcionQuery(ultimaTupla)
		rows, err := dbScrapping.Query(query)
		if err != nil {
			//TODO: Repetir en loop este error
			log.Println("Hay un error obteniendo las últimas ofertas\nError: ", err)
			time.Sleep(segundosDetenerError * time.Second)
			continue
		}
		log.Println("Se obtuvieron las últimas ofertas")
		defer rows.Close()
		//Se traducen las ofertas en base a la función del traductor
		log.Println("Traduciendo las últimas ofertas")
		ofertas, err := funcionTraductor(dbScrapping, rows)
		if err != nil {
			//TODO: Repetir en loop este error
			log.Println("Hay un error con la traducción de las consultas\nError: ", err)
			time.Sleep(segundosDetenerError * time.Second)
			continue
		}
		if len(ofertas) == 0 {
			log.Println("No hay ofertas nuevas o están vacías.")
			time.Sleep(segundosDetenerError * time.Second)
			continue
		}
		log.Printf("Se tradujeron las ofertas al struct")
		log.Printf("%+v\n", ofertas)
		for _, oferta := range ofertas {
			if oferta.IdComuna == 0 {
				//¿Existe la región de la oferta en la base de datos?
				var idRegion int
				log.Printf("Buscando si la región de '%s' la oferta ya existe...", oferta.Region)
				queryExisteRegion := fmt.Sprintf("select r.id from regiones r where lower(r.nombre) = '%s'", strings.ToLower(oferta.Region))

				filas, err := dbApp.Query(queryExisteRegion)
				empty := !filas.Next()
				filas.Scan(&idRegion)
				filas.Close()
				if err != nil {
					//TODO: Repetir en loop este error
					log.Println("Hay un error evaluando si la región de la oferta ya existe\nError: ", err)
					time.Sleep(segundosDetenerError * time.Second)
					continue
				}
				//Si no existe la región, se crea.
				if empty {
					log.Printf("No existe, creando la región '%s'\n", oferta.Region)
					queryCrearRegion := fmt.Sprintf("insert into regiones(nombre) values ('%s') returning id", oferta.Region)
					err = dbApp.QueryRow(queryCrearRegion).Scan(&idRegion)
					if err != nil {
						//TODO: Repetir en loop este error.
						log.Println("Hay un error creando la nueva región\nError: ", err)
						time.Sleep(segundosDetenerError * time.Second)
						continue
					}
					log.Printf("La región '%s' se creó existosamente.\n", oferta.Region)
				} else {
					log.Printf("Existe! y tiene un id %d\n", idRegion)
				}
				//¿Existe la comuna de la oferta en la base de datos?

				var idComuna int
				log.Printf("Buscando si la comuna '%s' de la oferta ya existe...\n", oferta.Comuna)
				queryExisteComuna := fmt.Sprintf("select c.id from comunas c where lower(c.nombre) = '%s'", strings.ToLower(oferta.Comuna))
				filas, err = dbApp.Query(queryExisteComuna)
				if err != nil {
					//TODO: Repetir en loop este error.
					log.Println("Hay un error evaluando si la comuna de la oferta ya existe\nError: ", err)
					time.Sleep(segundosDetenerError * time.Second)
					continue
				}
				empty = !filas.Next()
				filas.Scan(&idComuna)
				filas.Close()
				//Si no existe la comuna, se crea.
				if empty {
					log.Printf("No existe, creando la comuna '%s'\n", oferta.Comuna)
					queryCrearComuna := fmt.Sprintf("insert into comunas(nombre, id_region) values ('%s', %d) returning id", oferta.Comuna, idRegion)
					err = dbApp.QueryRow(queryCrearComuna).Scan(&idComuna)
					if err != nil {
						//TODO: Repetir en loop este error.
						log.Println("Hay un error creando la nueva comuna\nError: ", err)
						time.Sleep(segundosDetenerError * time.Second)
						continue
					}
					log.Printf("La comuna '%s' se creó existosamente.\n", oferta.Comuna)
				} else {
					log.Printf("Existe! y tiene un id %d\n", idComuna)
				}
				oferta.IdComuna = idComuna
			}

			//¿Existe el proveedor de la oferta en la base de datos?
			if oferta.IdProveedor == 0 {
				var idProveedor int
				log.Printf("Buscando si el proveedor '%s' de la oferta ya existe...\n", oferta.Proveedor)
				queryExisteProveedor := fmt.Sprintf("select p.id from proveedores p where lower(p.nombre) = '%s'", strings.ToLower(oferta.Proveedor))
				filas, err := dbApp.Query(queryExisteProveedor)
				if err != nil {
					//TODO: Repetir en loop este error.
					log.Println("Hay un error evaluando si el proveedor de la oferta ya existe\nError: ", err)
					time.Sleep(segundosDetenerError * time.Second)
					continue
				}
				empty := !filas.Next()
				filas.Scan(&idProveedor)
				filas.Close()
				//Si no existe el proveedor, se crea
				if empty {
					log.Printf("No existe, creando el proveedor '%s'\n", oferta.Proveedor)
					queryCrearProveedor := fmt.Sprintf("insert into proveedores(nombre) values ('%s') returning id", oferta.Proveedor)
					err = dbApp.QueryRow(queryCrearProveedor).Scan(&idProveedor)
					if err != nil {
						//TODO: Repetir en loop este error.
						log.Fatal("Hay un error creando el nuevo proveedor\nError: ", err)
						time.Sleep(segundosDetenerError * time.Second)
						continue
					}
					log.Printf("El proveedor '%s' se creó existosamente.\n", oferta.Proveedor)
				} else {
					log.Printf("Existe! y tiene un id %d\n", idProveedor)
				}
				oferta.IdProveedor = idProveedor
			}
			//Guardar nuevas ofertas en la db
			log.Println("Creando la nueva oferta turística...")
			queryCrearOferta := funcionCreacion(oferta)
			var idNuevaOferta int
			err = dbApp.QueryRow(queryCrearOferta).Scan(&idNuevaOferta)
			if err != nil {
				//TODO: Repetir en loop este error.
				log.Println("Hubo un error creando la nueva oferta turística\nError: ", err)
				//--------------------------Actualizar ID Scrapping-----------------------------------//
				ultimoId := ofertas[len(ofertas)-1].ID
				queryActualizarTupla := queryActualizarScrapping(ultimoId, nombreDBScrapping)
				err = dbApp.QueryRow(queryActualizarTupla).Scan(&auxScan)
				if err != nil {
					log.Println("Hubo un error actualizando el valor de la última tupla\nError: ", err)
					time.Sleep(segundosDetenerError * time.Second)
					continue
				}
				//------------------------------------------------------------------------------------//
				time.Sleep(segundosDetenerError * time.Second)
				continue
			}
			queryAsignarConsideracion := fmt.Sprintf("INSERT INTO oferta_consideraciones (id_oferta, id_consideracion) VALUES (%d, 1) returning id", idNuevaOferta)
			err = dbApp.QueryRow(queryAsignarConsideracion).Scan(&auxScan)
			if err != nil {
				log.Println("Ha habido un error relacionando la oferta con la consideración: ", err)
				time.Sleep(segundosDetenerError * time.Second)
				//--------------------------Actualizar ID Scrapping-----------------------------------//
				ultimoId := ofertas[len(ofertas)-1].ID
				queryActualizarTupla := queryActualizarScrapping(ultimoId, nombreDBScrapping)
				err = dbApp.QueryRow(queryActualizarTupla).Scan(&auxScan)
				if err != nil {
					log.Println("Hubo un error actualizando el valor de la última tupla\nError: ", err)
					time.Sleep(segundosDetenerError * time.Second)
					continue
				}
				//------------------------------------------------------------------------------------//
				continue
			}
			log.Println("Se creó una nueva oferta turística!")
			log.Printf("%+v\n", oferta)
		}
		log.Println("Se han actualizado todas las tuplas nuevas!")
		log.Println("Actualizando el último valor de la tupla")
		//Actualizar valor de la ultima tupla al último id de las ofertas
		//--------------------------Actualizar ID Scrapping-----------------------------------//
		ultimoId := ofertas[len(ofertas)-1].ID
		queryActualizarTupla := queryActualizarScrapping(ultimoId, nombreDBScrapping)
		err = dbApp.QueryRow(queryActualizarTupla).Scan(&auxScan)
		if err != nil {
			log.Println("Hubo un error actualizando el valor de la última tupla\nError: ", err)
			time.Sleep(segundosDetenerError * time.Second)
			continue
		}
		log.Println("Se ha actualizado el valor de la última tupla")
		//------------------------------------------------------------------------------------//
		time.Sleep(segundosDetenerFinal * time.Second)
	}
}

//-----------------------------------Query genérica para pruebas----------------------------------//
func QueryScrappingPrueba(ultimaTupla int) string {
	query := fmt.Sprintf("select o.id, o.comuna, o.fecha_inicio, o.fecha_final, o.nombre, o.precio, o.proveedor, o.region "+
		"from ofertas o "+
		"where o.id > %d"+
		"order by id", ultimaTupla)
	return query
}

func QueryScrappingPruebaCreacion(oferta modelos.OfertaTuristicaScrapping) string {
	queryCrearOferta := fmt.Sprintf("insert into ofertas_turisticas (nombre, precio, fecha_inicio, fecha_final, id_proveedor, id_comuna) "+
		"values ('%s', '%s', '%s', '%s', %d, %d) returning id",
		oferta.Nombre, oferta.Precio.String, oferta.FechaInicio, oferta.FechaFinal.String, oferta.IdProveedor, oferta.IdComuna)
	return queryCrearOferta
}

//------------------------------------------------------------------------------------------------//

//-------------------------------Query Chilepass--------------------------------------------------//
func QueryScrappingChilepass(ultimaTupla int) string {
	query := fmt.Sprintf("select e.id, e.title, d.id, now(), e.price, 1 as proveedor, e.ubication "+
		"from elements e, destinos d "+
		"where e.id > %d and e.destino_id = d.id "+
		"order by e.id", ultimaTupla)
	return query
}

func QueryScrappingChilepassCreacion(oferta modelos.OfertaTuristicaScrapping) string {
	queryCrearOferta := fmt.Sprintf("insert into ofertas_turisticas (nombre, precio, fecha_inicio, id_proveedor, id_comuna, ubicacion) "+
		"values ('%s', '%s', '%s', %d, %d, '%s') returning id",
		oferta.Nombre, oferta.Precio.String, oferta.FechaInicio, oferta.IdProveedor, oferta.IdComuna, oferta.Ubicacion.String)
	return queryCrearOferta
}

//------------------------------------------------------------------------------------------------//
