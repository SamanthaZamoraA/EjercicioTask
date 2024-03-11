package main

import (
	"fmt"
	"net/http"

	"github.com/Taks/internal/handler"
	"github.com/go-chi/chi"
)

func main() {

	//Inicializar las dependencias

	// Dependencia para el handler de tareas
	handler := handler.NewTaskHandler(nil, 0)

	//Dependencia para el router
	router := chi.NewRouter()

	//Registrar los endpoints
	router.Post("/post", handler.CreateTask())

	//Iniciar el servidor
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		return
	}
}
