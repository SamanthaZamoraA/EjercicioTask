package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Taks/internal"
	"github.com/Taks/internal/tools"
	"github.com/Taks/pkg/response"
)

// Se llama a la interfaz de servicio
type TaskHandler struct {
	sv internal.TaskService
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
func NewTaskHandler(sv internal.TaskService) *TaskHandler {
	//Se retorna el handler que contiene el servicio
	return &TaskHandler{
		sv: sv,
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

		//Paso 0: Leer el body
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.ResponseJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid request body"})
			return
		}

		// Paso 1: Decodificar el body y crear un map[string]any
		bodyMap := map[string]any{}
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": "invalid request body",
			})
			return
		}

		//Pso 2: Validar que todos los campos esten completos
		if err := tools.CheckFieldExistance(bodyMap, "tittle", "description", "done"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.ResponseJSON(w, http.StatusBadRequest, map[string]any{
					"message": fmt.Sprintf("%s is required", fieldError.Field),
				})
				return
			}

			response.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "internal server error",
			})
			return
		}

		// Paso 3: Crear una instancia de TaskRequest para almacenar los datos del cuerpo JSON
		var body TaskRequest

		// Paso 4: Decodificar el JSON del cuerpo de la solicitud y asignarlo a la estructura TaskRequest
		if err := json.Unmarshal(bytes, &body); err != nil {

			//Si ocurre un error decodificando el JSON, enviar una respuesta HTTP con error
			response.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"--Message": "invalid request body",
			})
			return
		}

		//process
		// Paso 5: Crear una instancia de Task a partir de los datos recibidos
		task := internal.Task{
			Tittle:      body.Tittle,
			Description: body.Description,
			Done:        body.Done,
		}

		// Paso 6: Agregar la tarea al mapa de tareas, usando el metodo Save del repositorio
		// Al Save se le pasa la tarea con los datos recibidos y en el repository se gestiona el guarda en el mapa y el id
		if err := t.sv.Save(&task); err != nil {

			// Se gestiona que tipo de error se produce y se envia la respuesta correspondiente
			switch {
			case errors.Is(err, internal.ErrTaskDuplicated):
				response.ResponseJSON(w, http.StatusConflict, map[string]any{"message": "task already exists"})
			case errors.Is(err, internal.ErrTaskInvalidField):
				response.ResponseJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid field"})
			default:
				response.ResponseJSON(w, http.StatusInternalServerError, map[string]any{"message": "internal server error"})
			}
			return
		}

		//response

		// Paso 7: Crear  una tarea en formatoJSON que se va a enviar como respuesta del handler
		data := TaskJSON{
			ID:          task.ID,
			Tittle:      task.Tittle,
			Description: task.Description,
			Done:        task.Done,
		}

		// Paso 8: Enviar una respuesta HTTP exitosa (201 Created) junto con los datos de la tarea creada
		response.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "task created successfully",
			"data":    data,
		})
	}
}
