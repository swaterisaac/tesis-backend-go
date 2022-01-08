package conexion

import (
	"fmt"
	"os"
)

func ObtenerConexionScrappingPrueba() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		os.Getenv("SCRAP_PRUEBA_HOST"),
		os.Getenv("SCRAP_PRUEBA_PORT"),
		os.Getenv("SCRAP_PRUEBA_USER"),
		os.Getenv("SCRAP_PRUEBA_PASSWORD"),
		os.Getenv("SCRAP_PRUEBA_NAME"),
		isSSL(),
	)
}
