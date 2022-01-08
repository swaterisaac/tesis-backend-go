package conexion

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

func isSSL() string {
	modoSSL := os.Getenv("SSL")
	esSSL, err := strconv.ParseBool(modoSSL)
	if err == nil && !esSSL {
		return "disable"
	} else {
		return "require"
	}
}

//Inicialización para cargar archivo .env (Variables de entorno)
func ObtenerEntornoConexion() string {
	fmt.Println(isSSL())
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		os.Getenv("PGDB_HOST"),
		os.Getenv("PGDB_PORT"),
		os.Getenv("PGDB_USER"),
		os.Getenv("PGDB_PASSWORD"),
		os.Getenv("PGDB_NAME"),
		isSSL(),
	)
}

func ConexionDb(psqlInfo string) (*sql.DB, error) {
	//--Conexión a la base de datos--//
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
