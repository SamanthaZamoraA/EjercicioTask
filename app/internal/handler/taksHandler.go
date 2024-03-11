package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Taks/internal"
	"github.com/Taks/internal/tools"
	"github.com/Taks/pkg/request"
	"github.com/Taks/pkg/response"
)

// Se crea una estructura de mapas para almacenar las tareas
type TaskHandler struct {
	tasks  map[int]internal.Task
	lastId int
}

// Se crea una estructura para almacenar las tareas en forma de requests
type TaskRequest struct {
	Tittle      string `json:"tittle"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Se crea una estructura para almacenar las tareas en forma de JSON
type TaskJSON struct {
	ID          int    `json:"id"`
	Tittle      string `json:"tittle"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Funcion para inicializar el handler de tareas
func NewTaskHandler(tasks map[int]internal.Task, lastId int) *TaskHandler {

	//Setear valores por defecto
	defaultTasks := make(map[int]internal.Task)
	defaultLastId := 0

	//Si es diferente de nil, setear los valores
	if tasks != nil {
		defaultTasks = tasks
	}

	if lastId != 0 {
		defaultLastId = lastId
	}

	//Retornar los valores por default
	return &TaskHandler{
		defaultTasks,
		defaultLastId,
	}
}

// --------------------- HANDLER DE CREATE ---------------------

// Metodo para crear una nueva tarea
func (t *TaskHandler) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request

		/*
			//validación de token
			token := r.Header.Get("Authorization")
			if token != "123" {
				response.ResponseJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"message": "unauthorized",
				})
				return
			}
		*/

		//Paso 0: Validar que todos los campos esten completos
		bodyMap := make(map[string]interface{})
		if err := request.RequestJSON(r, &bodyMap); err != nil {
			response.ResponseJSON(w, http.StatusBadRequest, map[string]interface{}{
				"message": "invalid request body",
			})
			return
		}

		if err := tools.CheckFieldExistance(bodyMap, "tittle", "description", "done"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.ResponseJSON(w, http.StatusBadRequest, map[string]interface{}{
					"message": fmt.Sprintf("%s is required", fieldError.Field),
				})
				return
			}

			response.ResponseJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal server error",
			})
			return
		}

		// Paso 1: Crear una instancia de TaskRequest para almacenar los datos del cuerpo JSON
		var body TaskRequest

		// Paso 2: Decodificar el JSON del cuerpo de la solicitud y asignarlo a la estructura TaskRequest
		if err := request.RequestJSON(r, &body); err != nil {

			//Si ocurre un error decodificando el JSON, enviar una respuesta HTTP con error
			response.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"Message": "invalid request body",
			})
			return
		}

		//process
		// Paso 3: Incrementar el ID de la tarea
		t.lastId++

		// Paso 4: Crear una instancia de Task a partir de los datos recibidos
		task := internal.Task{
			ID:          t.lastId,
			Tittle:      body.Tittle,
			Description: body.Description,
			Done:        body.Done,
		}

		// Paso 5: Agregar la tarea al mapa de tareas
		t.tasks[task.ID] = task

		//response

		// Paso 6: Crear  una tarea en formatoJSON
		data := TaskJSON{
			ID:          task.ID,
			Tittle:      task.Tittle,
			Description: task.Description,
			Done:        task.Done,
		}

		// Paso 7: Enviar una respuesta HTTP exitosa (201 Created) junto con los datos de la tarea creada
		response.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "task created successfully",
			"data":    data,
		})
	}
}
