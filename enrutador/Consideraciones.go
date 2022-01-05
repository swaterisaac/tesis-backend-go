package enrutador

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"tesis/busquedas"

	"github.com/gin-gonic/gin"
)

func ObtenerConsideracionesMedicas(c *gin.Context, db *sql.DB) {
	consideraciones, err, codigo := busquedas.ObtenerConsideracionesMedicas(db)

	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "No hay consideraciones médicas",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Hay un error con el servidor",
		})
		return
	}
	consideraciones = consideraciones[1:]
	c.JSON(http.StatusOK, consideraciones)
}

func ObtenerConsideracionesPorOferta(c *gin.Context, db *sql.DB) {
	idOferta := c.Query("idOferta")
	idOfertaNum, err := strconv.Atoi(idOferta)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		log.Println("Error: ", err)
		return
	}
	consideraciones, err := busquedas.ObtenerConsideracionesPorOferta(db, idOfertaNum)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		log.Println("Error: ", err)
		return
	}
	c.JSON(http.StatusOK, consideraciones)
}
