package main

import (
	"fmt"
	"os"
	"tesis/conexion"
	"tesis/enrutador"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("No se encontró el .env", r)
		}
	}()
	err := godotenv.Load(".env.dev")
	if err != nil {
		panic(err)
	}
}
func main() {

	//Puerto de heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	//---------------------Conexión base de datos principal de la aplicación--------------//
	db, err := conexion.ConexionDb(conexion.ObtenerEntornoConexion())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//-----------------------------------------------------------------------------------//
	//---------------------Conexión base de datos de Scraping Prueba---------------------//

	dbScrappingPrueba, err := conexion.ConexionDb(conexion.ObtenerConexionScrappingPrueba())
	if err != nil {
		panic(err)
	}
	defer dbScrappingPrueba.Close()

	//Listener para scrapping
	//go conexion.ListenerGeneral(db, dbScrappingPrueba, conexion.QueryScrappingPrueba, conexion.TraductorScrappingGeneral, conexion.QueryScrappingPruebaCreacion, "prueba")
	go conexion.ListenerGeneral(db, dbScrappingPrueba, conexion.QueryScrappingChilepass, conexion.TraductorScrappingChilepass, conexion.QueryScrappingChilepassCreacion, "chilepass")

	//elasticsearch
	cloudID := os.Getenv("ELASTICSEARCH_CLIENT_ID")
	apiKey := os.Getenv("ELASTICSEARCH_API_KEY")

	es, err := elasticsearch.NewDefaultClient()
	if cloudID != "" && apiKey != "" {
		cfg := elasticsearch.Config{
			CloudID: cloudID,
			APIKey:  apiKey,
		}
		es, err = elasticsearch.NewClient(cfg)
	}

	if err != nil {
		panic(err)
	}
	//Builds vacíos para Heroku: 2
	/*
		//Para debugear elasticsearch, descomentar esto
		modelos, err := busquedas.ObtenerOfertasRecomendacion(db, es, 3, 6, 1)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v", modelos)
	*/
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.GET("/ofertas", func(c *gin.Context) {
		enrutador.ObtenerOfertasUsuario(c, db, es)
	})
	router.GET("/ofertasQuery", func(c *gin.Context) {
		enrutador.ObtenerOfertasQuery(c, db, es)
	})
	router.GET("/usuarioPorCorreo", func(c *gin.Context) {
		enrutador.ObtenerUsuarioPorCorreo(c, db)
	})
	router.GET("/obtenerRegiones", func(c *gin.Context) {
		enrutador.ObtenerRegiones(c, db)
	})
	router.GET("/obtenerComunas", func(c *gin.Context) {
		enrutador.ObtenerComunasPorRegion(c, db)
	})
	router.GET("/obtenerProveedores", func(c *gin.Context) {
		enrutador.ObtenerProveedores(c, db)
	})
	router.GET("/obtenerConsideraciones", func(c *gin.Context) {
		enrutador.ObtenerConsideracionesMedicas(c, db)
	})
	router.GET("/obtenerOfertaPorId", func(c *gin.Context) {
		enrutador.ObtenerOfertaPorId(c, db)
	})
	router.GET("/obtenerConsideracionesPorOferta", func(c *gin.Context) {
		enrutador.ObtenerConsideracionesPorOferta(c, db)
	})
	router.POST("/crearUsuario", func(c *gin.Context) {
		enrutador.CrearUsuario(c, db)
	})
	router.POST("/crearHistorialBusqueda", func(c *gin.Context) {
		enrutador.CrearHistorialBusqueda(c, db)
	})
	router.POST("/crearHistoriales", func(c *gin.Context) {
		enrutador.CrearHistoriales(c, db)
	})
	router.PUT("/editarUsuario", func(c *gin.Context) {
		enrutador.EditarUsuario(c, db)
	})
	router.Run(":" + port)
}
