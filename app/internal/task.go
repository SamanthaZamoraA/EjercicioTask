package internal

import "errors"

/*
	Este archivo es el domain que debe ir en la raiz de internal

	La comunicación es la siguiente:

	1. El handler llama al service
	2. El service llama al repository


	Se declara la estructura base - > Task
	Se declaran errores genericos personalizados, que las interfaces puedan retornar

	Se declara la interface de repositorio - > TaskRepository
	Se declara la interface de servicio - > TaskService

*/

// Creación de la estructura Task
type Task struct {
	ID          int
	Tittle      string
	Description string
	Done        bool
}

var (
	// Error para cuando no se encuentra la tarea
	ErrTaskNotFound = errors.New("task not found")

	//Erro para cuando esta duplicada la tarea
	ErrTaskDuplicated = errors.New("task duplicated")

	//Error al procesar la tarea
	ErrTaskProcessing = errors.New("task processing")

	//Error interno
	ErrTaskInternal = errors.New("task internal error")

	//Error campo invalido
	ErrTaskInvalidField = errors.New("task invalid field")

	//Error en el service
	ErrTaskService = errors.New("task service can´t be processed")
)

// Interfaz de repository
type TaskRepository interface {

	//Debe ser un puntero de Task porque se va a trabajar con el ultimo ID
	Save(task *Task) (err error)

	//Actualizar y sino esta devuelve error
	Update(task Task) (err error)

	//Actualizar parcialmente
	UpdatePartial(id int, fields map[string]any) (err error)

	//Eliminar una tarea
	Delete(id int) (err error)

	//Obtener todas las tareas

	//Obtener por id
	GetByID(id int) (task Task, err error)
}

// Interfaz de service
type TaskService interface {
	Save(task *Task) (err error)

	Update(task Task) (err error)

	UpdatePartial(id int, fields map[string]any) (err error)

	Delete(id int) (err error)

	GetByID(id int) (task Task, err error)
}
