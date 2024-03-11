package service

import (
	"github.com/Taks/internal"
)

/*
Estructura de TaskService que se conecta con el repositorio.
Desde el servicio se acceder a todo lo que ofrece el repositorio
*/
type TaskService struct {
	repository internal.TaskRepository
}

// Funcion para inicializar el servicio de tareas
func NewTaskService(rp internal.TaskRepository) *TaskService {
	return &TaskService{
		repository: rp,
	}
}

// Funcion para implementar el metodo Save de la interfaz TaskService
func (t *TaskService) Save(task *internal.Task) (err error) {
	err = t.repository.Save(task)
	return
}
