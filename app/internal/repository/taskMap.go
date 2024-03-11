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

// Funcion para crear una tareas
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

// Funcion para actualizar una tarea
func (t *TaskMap) Update(task internal.Task) (err error) {
	//Verificar que exista
	if _, ok := (*t).db[(task).ID]; !ok {
		err = internal.ErrTaskNotFound
		return
	}

	//Verificar que no exista otra tarea con el mismo titulo
	for _, t := range (*t).db {
		if t.ID != task.ID && t.Tittle == task.Tittle {
			err = internal.ErrTaskDuplicated
			return
		}
	}

	//Actualizar la tarea
	(*t).db[(task).ID] = task
	return
}

// Funcion para actualizar parcialmente una tarea
func (t *TaskMap) UpdatePartial(id int, fields map[string]any) (err error) {
	//Verificar que exista
	task, ok := (*t).db[id]
	if !ok {
		err = internal.ErrTaskNotFound
		return
	}

	//Actualizar la tarea
	for key, value := range fields {

		//Verificar que el campo sea valido
		switch key {
		case "tittle", "Tittle":

			//Verificar que el valor sea un string
			tittle, ok := value.(string)
			if !ok {
				err = internal.ErrTaskInvalidField
				return
			}

			// Verificar que no exista otra tarea con el mismo titulo
			for _, t := range (*t).db {
				if t.ID != id && t.Tittle == tittle {
					err = internal.ErrTaskDuplicated
					return
				}
			}

			//Actualizar el titulo
			task.Tittle = tittle
		case "description", "Description":
			task.Description, ok = value.(string)
			if !ok {
				err = internal.ErrTaskNotFound
				return
			}
		case "done", "Done":
			task.Done, ok = value.(bool)
			if !ok {
				err = internal.ErrTaskNotFound
				return
			}
		default:
		}
	}

	//Actualizar la tarea
	(*t).db[id] = task
	return
}

// Funcion para eliminar una tarea
func (t *TaskMap) Delete(id int) (err error) {
	// Validar que exista
	_, ok := (*t).db[id]
	if !ok {
		err = internal.ErrTaskNotFound
		return
	}

	// Eliminar la tarea
	delete((*t).db, id)
	return
}
