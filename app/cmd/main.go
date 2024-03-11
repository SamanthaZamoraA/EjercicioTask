package main

import (
	"fmt"
	"net/http"

	"github.com/Taks/internal/handler"
	"github.com/Taks/internal/repository"
	"github.com/Taks/internal/service"
	"github.com/go-chi/chi"
)

func main() {

	//Inicializar las dependencias

	//Dependencia del repository
	rp := repository.NewTaskMap(nil, 0)

	//Dependencia del service
	sv := service.NewTaskService(rp)

	//Dependencia del handler
	h := handler.NewTaskHandler(sv)

	//Dependencia para el router
	router := chi.NewRouter()

	//Registrar los endpoints
	router.Route("/task", func(r chi.Router) {
		//Método POST
		r.Post("/post", h.CreateTask())

		//Método PUT
		r.Put("/put/{id}", h.UpdateTask())

		//Método PATCH
		r.Patch("/patch/{id}", h.UpdatePartialTask())

		//Método DELETE
		r.Delete("/delete/{id}", h.DeleteTask())
	})

	//Iniciar el servidor
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		return
	}
}
