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
	log.Println("servidor corriendo....")
	http.ListenAndServe(":8080", nil)
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
