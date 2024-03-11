package repository

import "github.com/Taks/internal"

// Esta es una implementacion de la interfaz TaskRepository basada en un mapa
type TaskMap struct {
	db     map[int]internal.Task
	lastId int
}

// Funcion para inicializar el repositorio de tareas
func NewTaskMap(mapa map[int]internal.Task, lastId int) *TaskMap {

	//Setear valores por defecto
	defaultTasks := make(map[int]internal.Task)
	defaultLastId := 0

	//Si es diferente de nil, setear los valores
	if mapa != nil {
		defaultTasks = mapa
	}

	if lastId != 0 {
		defaultLastId = lastId
	}

	//Retornar el repositorio
	return &TaskMap{
		defaultTasks,
		defaultLastId,
	}
}

// Funcion para inicializar el repositorio de tareas
func (t *TaskMap) Save(task *internal.Task) (err error) {

	//Se valida que la tarea no este duplicada
	for _, value := range (*t).db {
		if value.Tittle == (*task).Tittle {
			err = internal.ErrTaskDuplicated
			return
		}
	}

	//Se incrementa el ultimo ID
	(*t).lastId++

	//Se asigna a la tarea el ultimo ID
	(*task).ID = (*t).lastId

	//Se guarda la tarea en el mapa
	(*t).db[(*task).ID] = *task

	return
}
