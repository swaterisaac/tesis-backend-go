package conexion

import (
	"database/sql"
	"fmt"
	"os"
)

//Inicialización para cargar archivo .env (Variables de entorno)
/*
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
*/
func ObtenerEntornoConexion() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		os.Getenv("PGDB_HOST"),
		os.Getenv("PGDB_PORT"),
		os.Getenv("PGDB_USER"),
		os.Getenv("PGDB_PASSWORD"),
		os.Getenv("PGDB_NAME"),
	)
}

func ConexionDb(psqlInfo string) (*sql.DB, error) {
	//--Conexión a la base de datos--//
	fmt.Println("AKI->", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
