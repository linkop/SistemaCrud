package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

//funcion para conectar con la base datos

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "sistemacrud"
)

func conxionBD() (conexion *sql.DB) {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	conexion, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	//defer conexion.Close()
	//err = conexion.Ping()
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("Ping OK")
	return conexion
} //fin de la conexion a la base de datos

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	//rutas
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("servidor corriendo....")
	http.ListenAndServe(":8080", nil)
}

// funcion para borrar datos de la tabla
func Borrar(rw http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conexionEstablecida := conxionBD()
	borrarRegistro, err := conexionEstablecida.Prepare(`delete from empleado where id=$1`)

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)

	http.Redirect(rw, r, "/", 301)
}

// funciona para leer los elementos de la tabla empleado
type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(rw http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conxionBD()
	Regristros, err := conexionEstablecida.Query(`select * from empleado`)
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for Regristros.Next() {
		var id int
		var nombre, correo string
		err = Regristros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	fmt.Println(arregloEmpleado)

	//fmt.Fprintf(rw, "Hola mundo")
	//con esto accedemos a la plantilla inicio
	plantillas.ExecuteTemplate(rw, "inicio", arregloEmpleado)
	//antes de mostrar los datos dinamicamente desde la base de datos
	//plantillas.ExecuteTemplate(rw, "inicio", nil)
}

// funcion que me permite editar elementos en la base de datos
func Editar(rw http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conexionEstablecida := conxionBD()
	Regristro, err := conexionEstablecida.Query(`select * from empleado where id=$1`, idEmpleado)

	empleado := Empleado{}
	for Regristro.Next() {
		var id int
		var nombre, correo string
		err = Regristro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(rw, "editar", empleado)
}

// funcion handler que nos envia a la plantilla que tiene
// el formulario para agregar usuario
func Crear(rw http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(rw, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conxionBD()
		insertarRegristro, err := conexionEstablecida.Prepare(`insert into empleado(nombre,correo) VALUES ($1,$2)`)

		if err != nil {
			panic(err.Error())
		}
		insertarRegristro.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)

	}
}

// funcion que actualiza los campos de la base de datos
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conxionBD()
		modificarRegristro, err := conexionEstablecida.Prepare(`update empleado set nombre=$1 ,correo=$2 where id=$3`)
		//update empleado set nombre='soyelsiete' ,correo='qq@gmail.com' where id=7
		if err != nil {
			panic(err.Error())
		}
		modificarRegristro.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)

	}
}
