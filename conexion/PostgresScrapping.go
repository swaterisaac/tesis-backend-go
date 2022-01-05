package conexion

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

//Inicialización para cargar archivo .env (Variables de entorno)
func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func ObtenerConexionScrappingPrueba() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("SCRAP_PRUEBA_HOST"),
		os.Getenv("SCRAP_PRUEBA_PORT"),
		os.Getenv("SCRAP_PRUEBA_USER"),
		os.Getenv("SCRAP_PRUEBA_PASSWORD"),
		os.Getenv("SCRAP_PRUEBA_NAME"),
	)
}
