package modelos

/*
Struct que representa a un usuario de la base de datos.

ID: int Con el id del usuario
Nombre: string con el nombre del usuario
Correo: string Con el correo del usuario
FechaNacimiento: string Con la fecha de nacimiento del usuario
IdComuna: int Con el id de la comuna a la que pertenece el usuario
*/
type Usuario struct {
	ID              int
	Nombre          string `validate:"required" form:"nombre"`
	Correo          string `validate:"required,email" form:"correo"`
	FechaNacimiento string `validate:"required" form:"fechaNacimiento"`
	IdComuna        int    `validate:"required" form:"idComuna"`
	Consideraciones []int  `form:"consideraciones"`
}

type UsuarioApp struct {
	ID                     int      `json:"id"`
	Nombre                 string   `json:"nombre"`
	Correo                 string   `json:"correo"`
	FechaNacimiento        string   `json:"fechaNacimiento"`
	Comuna                 string   `json:"comuna"`
	Region                 string   `json:"region"`
	ConsideracionesMedicas []string `json:"consideracionesMedicas"`
}

type UsuarioHistorial struct {
	Correo   string `form:"correo"`
	Busqueda string `form:"busqueda"`
}

type UsuarioHistoriales struct {
	Correo   string `form:"correo"`
	IdOferta int    `form:"idOferta"`
}
