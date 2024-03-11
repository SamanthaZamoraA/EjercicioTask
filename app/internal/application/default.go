package application

import (
	"fmt"
	"net/http"

	"github.com/Taks/internal/handler"
	"github.com/Taks/internal/repository"
	"github.com/Taks/internal/service"
	"github.com/go-chi/chi"
)

// Default  es una implenetación de application.
type Default struct {
	// addr es la dirección donde se va a ejecutar el servidor.
	addr string
}

// NewDefault retorns a new Default application.
func NewDefault(addr string) *Default {
	if addr == "" {
		addr = ":8080"
	}

	return &Default{
		addr: addr,
	}
}

// Método para correr con los paths
func (a *Default) Run() (err error) {
	// dependencias
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

		//Método GETBYID
		r.Get("/get/{id}", h.GetTaskByID())
	})

	// Iniciar el servidor
	if err := http.ListenAndServe(a.addr, router); err != nil {
		return fmt.Errorf("error al iniciar el servidor: %v", err)
	}

	return nil
}
