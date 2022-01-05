package enrutador

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"tesis/busquedas"

	"github.com/gin-gonic/gin"
)

func ObtenerRegiones(c *gin.Context, db *sql.DB) {
	regiones, err, codigo := busquedas.ObtenerRegiones(db)
	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No hay regiones",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Hay un error con el servidor",
		})
		return
	}
	fmt.Printf("%+v\n", regiones)
	c.JSON(http.StatusOK, regiones)
}

func ObtenerComunasPorRegion(c *gin.Context, db *sql.DB) {
	idRegion := c.Query("region")
	idRegionNum, err := strconv.Atoi(idRegion)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	comunas, err, codigo := busquedas.ObtenerComunasPorRegion(db, idRegionNum)
	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No existen comunas de esa region",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Hay un error con el servidor",
		})
		return
	}
	c.JSON(http.StatusOK, comunas)
}
