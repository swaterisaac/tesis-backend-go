package enrutador

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tesis/busquedas"
	"tesis/modelos"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/gin-gonic/gin"
)

func ObtenerOfertasUsuario(c *gin.Context, db *sql.DB, es *elasticsearch.Client) {
	correo := c.Query("correo")
	tamanio := c.DefaultQuery("tamanio", "6")
	pagina := c.DefaultQuery("pagina", "1")

	tamanioNum, errTamanio := strconv.Atoi(tamanio)
	paginaNum, errPagina := strconv.Atoi(pagina)

	if errTamanio != nil || errPagina != nil {
		log.Println("Error tamanio: ", errTamanio)
		log.Println("Error pagina: ", errPagina)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Tamaño y pagina deben ser valores numéricos",
		})
		return
	}
	if tamanioNum <= 0 || paginaNum <= 0 || tamanioNum > 1000 || paginaNum > 1000 {
		log.Println("Error: Tamaño y pagina deben ser positivos distintos de 0 y menores a 1000")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Tamaño y pagina deben ser positivos distintos de 0 y menores a 1000",
		})
		return
	}
	if correo == "" {
		log.Println("Error: No se ha ingresado correo")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "No se ha ingresado correo",
		})
		return
	}
	id, err, codigo := busquedas.ObtenerIdUsarioPorCorreo(db, correo)

	if codigo == http.StatusNotFound {
		log.Println("Error: ", err)
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "El correo no está registrado",
		})
		return
	}
	if err != nil {
		log.Println("Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error de servidor",
			"error": err,
		})
		return
	}
	ofertas, err := busquedas.ObtenerOfertasRecomendacion(db, es, id, tamanioNum, paginaNum)
	if err != nil {
		log.Println("Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	if len(ofertas) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No hay ofertas",
		})
		return
	}
	c.JSON(http.StatusOK, ofertas)

}

func ObtenerOfertasQuery(c *gin.Context, db *sql.DB, es *elasticsearch.Client) {
	correo := c.Query("correo")
	tamanio := c.DefaultQuery("tamanio", "6")
	pagina := c.DefaultQuery("pagina", "1")
	query := c.Query("query")
	filtroComuna := c.Query("filtroComuna")
	filtroRegion := c.Query("filtroRegion")
	filtroProveedor := c.Query("filtroProveedor")
	if filtroComuna != "" {
		query += " " + filtroComuna
	}
	if filtroRegion != "" {
		query += " " + filtroRegion
	}
	if filtroProveedor != "" {
		query += " " + filtroProveedor
	}

	existeFiltro := filtroComuna != "" || filtroRegion != "" || filtroProveedor != ""

	tamanioNum, errTamanio := strconv.Atoi(tamanio)
	paginaNum, errPagina := strconv.Atoi(pagina)

	if errTamanio != nil || errPagina != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Tamaño y pagina deben ser valores numéricos",
		})
		return
	}
	if tamanioNum <= 0 || paginaNum <= 0 || tamanioNum > 1000 || paginaNum > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Tamaño y pagina deben ser positivos distintos de 0 y menores a 1000",
		})
		return
	}
	if correo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "No se ha ingresado correo",
		})
		return
	}

	id, err, codigo := busquedas.ObtenerIdUsarioPorCorreo(db, correo)

	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "El correo no está registrado",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error de servidor",
			"error": err,
		})
		return
	}
	ofertas, err := busquedas.ObtenerOfertasQuery(db, es, id, query, tamanioNum, paginaNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	if len(ofertas) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No hay ofertas",
		})
		return
	}
	var ofertasFiltradas []modelos.OfertaTuristica
	if existeFiltro {
		fmt.Println("ALO")
		for _, oferta := range ofertas {
			if (oferta.Comuna == filtroComuna || filtroComuna == "") && (oferta.Region == filtroRegion || filtroRegion == "") && (oferta.Proveedor == filtroProveedor || filtroProveedor == "") {
				ofertasFiltradas = append(ofertasFiltradas, oferta)
			}
		}
		if len(ofertasFiltradas) == 0 {
			c.JSON(http.StatusNoContent, gin.H{
				"error": "No hay ofertas",
			})
			return
		}
		c.JSON(http.StatusOK, ofertasFiltradas)
		return
	}
	c.JSON(http.StatusOK, ofertas)
}

func ObtenerOfertaPorId(c *gin.Context, db *sql.DB) {
	idOferta := c.Query("idOferta")
	idOfertaNum, err := strconv.Atoi(idOferta)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"msg":   "No se ha ingresado idOferta como valor numérico",
		})
		return
	}
	oferta, err, codigo := busquedas.ObtenerOfertaPorId(db, idOfertaNum)

	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No existe oferta con ese id",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "Error de servidor",
			"error": err,
		})
		log.Println("Error: ", err)
		return
	}
	c.JSON(http.StatusOK, oferta)
}

func ObtenerProveedores(c *gin.Context, db *sql.DB) {
	proveedores, err := busquedas.ObtenerProveedores(db)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No hay proveedores",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, proveedores)
}
