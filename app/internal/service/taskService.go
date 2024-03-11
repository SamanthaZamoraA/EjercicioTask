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

// Funcion para implementar el metodo Update de la interfaz TaskService
func (t *TaskService) Update(task internal.Task) (err error) {
	err = t.repository.Update(task)
	return
}

// Funcion para implementar el metodo UpdatePartial de la interfaz TaskService
func (t *TaskService) UpdatePartial(id int, fields map[string]any) (err error) {
	err = t.repository.UpdatePartial(id, fields)
	return
}

// Funcion para implementar el metodo Delete de la interfaz TaskService
func (t *TaskService) Delete(id int) (err error) {
	err = t.repository.Delete(id)
	return
}
