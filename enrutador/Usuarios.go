package enrutador

import (
	"database/sql"
	"log"
	"net/http"
	"tesis/busquedas"
	"tesis/modelos"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

func ObtenerUsuarioPorCorreo(c *gin.Context, db *sql.DB) {
	correo := c.Query("correo")

	if correo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "No se ha ingresado correo",
		})
		return
	}

	usuario, err, codigo := busquedas.ObtenerUsuarioPorCorreo(db, correo)

	if codigo == http.StatusNotFound {
		c.JSON(http.StatusNoContent, gin.H{
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
	c.JSON(http.StatusOK, usuario)
}

func CrearUsuario(c *gin.Context, db *sql.DB) {
	var usuario modelos.Usuario
	if err := c.Bind(&usuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error con el traspaso de datos",
		})
		return
	}
	validate := validator.New()
	err := validate.Struct(&usuario)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if codigo := busquedas.CrearUsuario(db, usuario); codigo == http.StatusInternalServerError {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Hay un error con el servidor",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Se ha creado el usuario",
	})
}

func EditarUsuario(c *gin.Context, db *sql.DB) {
	var usuario modelos.Usuario
	if err := c.Bind(&usuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error con el traspaso de datos",
		})
		log.Println("Error: \n", err)
		return
	}

	usuarioID, err, _ := busquedas.ObtenerUsuarioPorCorreo(db, usuario.Correo)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error con el servidor",
		})
		log.Println("Error: \n", err)
		return
	}
	usuario.ID = usuarioID.ID

	err = busquedas.EditarUsuario(db, usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Error con el servidor",
		})
		log.Println("Error: \n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
