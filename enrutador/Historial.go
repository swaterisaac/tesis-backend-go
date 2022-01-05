package enrutador

import (
	"database/sql"
	"net/http"
	"tesis/busquedas"
	"tesis/modelos"

	"github.com/gin-gonic/gin"
)

func CrearHistorialBusqueda(c *gin.Context, db *sql.DB) {
	//Este modelo está dentro de usuario.
	var datos modelos.UsuarioHistorial

	err := c.Bind(&datos)
	if err != nil || datos.Busqueda == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
			"msg":   "Datos entregados incorrectos",
		})
		return
	}
	usuario, err, codigo := busquedas.ObtenerUsuarioPorCorreo(db, datos.Correo)
	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
			"msg":   "No existe usuario con ese correo",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error del servidor",
		})
		return
	}

	codigo = busquedas.CrearHistorialBusqueda(db, usuario.ID, datos.Busqueda)

	if codigo == http.StatusInternalServerError {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error del servidor",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"msg": "Se ha creado el historial de busqueda",
	})
}

func CrearHistoriales(c *gin.Context, db *sql.DB) {
	//Este modelo está dentro de usuario.
	var datos modelos.UsuarioHistoriales

	err := c.Bind(&datos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
			"msg":   "Datos entregados incorrectos",
		})
		return
	}
	usuario, err, codigo := busquedas.ObtenerUsuarioPorCorreo(db, datos.Correo)
	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
			"msg":   "No existe usuario con ese correo",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error del servidor",
		})
		return
	}

	codigo = busquedas.CrearHistoriales(db, usuario.ID, datos.IdOferta)
	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "No existe esa oferta turística",
		})
		return
	}
	if codigo != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error del servidor",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"msg": "Se ha creado los historiales",
	})
}
