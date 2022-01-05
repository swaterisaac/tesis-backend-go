package busquedas

import (
	"database/sql"
	"fmt"
	"net/http"
	"tesis/modelos"
)

/*
Entradas:
db: *sql.DB sacado de la librería database/sql. Manejador de conexión a bases de datos mysql
(pgsql en este caso)

Salidas:
[]modelos.Usuario: sacado de un struct en el paquete modelos que representa al usuario y sus datos
error: Error en caso de que algo salga mal

Descripción:
Retorna todos los usuarios de la DB en un array de struct de Usuario.
*/
func ObtenerUsuarios(db *sql.DB) ([]modelos.Usuario, error) {
	rows, err := db.Query("SELECT * FROM usuarios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var usuarios []modelos.Usuario

	for rows.Next() {
		var usuario modelos.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nombre, &usuario.Correo,
			&usuario.FechaNacimiento, &usuario.IdComuna); err != nil {
			return usuarios, err
		}
		usuarios = append(usuarios, usuario)
	}
	if err = rows.Err(); err != nil {
		return usuarios, err
	}
	return usuarios, nil
}

/*
Entradas:
db: *sql.DB sacado de la librería database/sql. Manejador de conexión a bases de datos mysql
(pgsql en este caso)
idUsuario: int que indica el id del usuario a sacar los datos.

Salidas:
string: Comuna a la que pertenece un usuario
string: Región a la que pertenece un usuario
error: Error en caso de que algo salga mal

Descripción:
En base a un usuario específico, retorna la comuna y región a la que pertenece.
*/
func ObtenerUbicacionUsuario(db *sql.DB, idUsuario int) (string, string, error) {
	query := fmt.Sprintf("SELECT r.nombre, c.nombre "+
		"FROM usuarios u, regiones r, comunas c "+
		"WHERE u.id = %d and c.id = u.id_comuna and r.id = c.id_region", idUsuario)
	rows, err := db.Query(query)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	comuna, region := "", ""
	for rows.Next() {
		if err := rows.Scan(&region, &comuna); err != nil {
			return comuna, region, err
		}
	}
	if err = rows.Err(); err != nil {
		return comuna, region, err
	}

	return comuna, region, nil
}

/*
Entradas:
db: *sql.DB sacado de la librería database/sql. Manejador de conexión a bases de datos mysql
(pgsql en este caso)
idUsuario: int que indica el id del usuario a sacar los datos.

Salidas:
[]string: Lista de string con las consideraciones médicas del usuario
error: Error en caso de que algo salga mal

Descripción:
En base a un usuario específico retorna todas las consideraciones médicas asociadas.
*/
func ObtenerConsideracionesMedicasUsuario(db *sql.DB, idUsuario int) ([]string, error) {
	query := fmt.Sprintf("SELECT cm.nombre "+
		"FROM usuarios u, consideraciones_medicas cm , usuario_consideraciones uc "+
		"WHERE u.id = %d and uc.id_usuario = u.id and uc.id_consideracion = cm.id", idUsuario)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var consideraciones []string

	for rows.Next() {
		var consideracion string
		if err := rows.Scan(&consideracion); err != nil {
			return consideraciones, err
		}
		consideraciones = append(consideraciones, consideracion)
	}

	if err := rows.Err(); err != nil {
		return consideraciones, err
	}

	return consideraciones, nil
}

func ObtenerIdUsarioPorCorreo(db *sql.DB, correo string) (int, error, int) {
	query := fmt.Sprintf("SELECT id FROM usuarios u WHERE u.correo = '%s'", correo)
	var id int
	err := db.QueryRow(query).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, err, 404
	}
	if err != nil {
		return 0, err, 500
	}
	return id, nil, 200
}

func ObtenerUsuarioPorCorreo(db *sql.DB, correo string) (modelos.UsuarioApp, error, int) {
	query := fmt.Sprintf("SELECT u.id, u.nombre, u.correo, u.fecha_nacimiento, c.nombre, r.nombre "+
		"FROM usuarios u, comunas c, regiones r "+
		"WHERE u.correo = '%s' AND u.id_comuna = c.id AND c.id_region = r.id", correo)
	var usuario modelos.UsuarioApp
	err := db.QueryRow(query).Scan(&usuario.ID, &usuario.Nombre, &usuario.Correo, &usuario.FechaNacimiento, &usuario.Comuna, &usuario.Region)

	if err == sql.ErrNoRows {
		return usuario, err, 404
	}
	if err != nil {
		return usuario, err, 500
	}
	queryConsideraciones := fmt.Sprintf("SELECT cm.nombre "+
		"FROM usuarios u, consideraciones_medicas cm, usuario_consideraciones uc "+
		"WHERE u.id = %d AND u.id = uc.id_usuario AND cm.id = uc.id_consideracion", usuario.ID)

	rows, err := db.Query(queryConsideraciones)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return usuario, nil, 200
	}
	for rows.Next() {
		var consideracion string
		if err := rows.Scan(&consideracion); err != nil {
			return usuario, err, 500
		}
		usuario.ConsideracionesMedicas = append(usuario.ConsideracionesMedicas, consideracion)
	}
	return usuario, nil, 200
}

func CrearUsuario(db *sql.DB, usuario modelos.Usuario) int {
	queryUsuario := fmt.Sprintf("INSERT INTO usuarios (nombre, correo, fecha_nacimiento, id_comuna) "+
		"VALUES ('%s', '%s', '%s', %d) RETURNING id", usuario.Nombre, usuario.Correo, usuario.FechaNacimiento, usuario.IdComuna)
	var idUsuario int

	if err := db.QueryRow(queryUsuario).Scan(&idUsuario); err != nil {
		return http.StatusInternalServerError
	}
	for _, idConsideracion := range usuario.Consideraciones {
		queryConsideraciones := fmt.Sprintf("INSERT INTO usuario_consideraciones (id_usuario, id_consideracion) "+
			"VALUES (%d, %d)", idUsuario, idConsideracion)
		if err := db.QueryRow(queryConsideraciones).Err(); err != nil {
			return http.StatusInternalServerError
		}
	}
	return http.StatusCreated
}

func agregarConsideracionUsuario(db *sql.DB, idUsuario int, idConsideracion int) error {
	query := fmt.Sprintf("INSERT INTO usuario_consideraciones (id_usuario, id_consideracion) "+
		"VALUES (%d, %d)", idUsuario, idConsideracion)
	err := db.QueryRow(query).Err()
	return err
}

func borrarConsideracionesUsuario(db *sql.DB, idUsuario int) error {
	query := fmt.Sprintf("DELETE FROM usuario_consideraciones WHERE id_usuario = %d", idUsuario)
	err := db.QueryRow(query).Err()
	return err
}

func EditarUsuario(db *sql.DB, usuario modelos.Usuario) error {
	camposActualizar := ""
	primerCampo := false
	if usuario.FechaNacimiento != "" {
		camposActualizar += fmt.Sprintf("fecha_nacimiento = '%s'", usuario.FechaNacimiento)
		primerCampo = true
	}
	if usuario.Nombre != "" {
		if primerCampo {
			camposActualizar += ", "
		}
		camposActualizar += fmt.Sprintf("nombre = '%s'", usuario.Nombre)
		primerCampo = true
	}
	if usuario.IdComuna != 0 {
		if primerCampo {
			camposActualizar += ", "
		}
		camposActualizar += fmt.Sprintf("id_comuna = %d", usuario.IdComuna)
	}
	if camposActualizar != "" {
		fmt.Println("AQUI: ", camposActualizar)
		queryActualizar := fmt.Sprintf("UPDATE usuarios SET %s WHERE id = %d", camposActualizar, usuario.ID)
		err := db.QueryRow(queryActualizar).Err()
		if err != nil {
			return err
		}
	}
	if len(usuario.Consideraciones) > 0 {
		err := borrarConsideracionesUsuario(db, usuario.ID)
		if err != nil {
			return err
		}
	}
	for _, idConsideracion := range usuario.Consideraciones {
		if err := agregarConsideracionUsuario(db, usuario.ID, idConsideracion); err != nil {
			return err
		}
	}
	return nil
}
