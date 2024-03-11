package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Taks/internal"
	"github.com/Taks/internal/tools"
	"github.com/Taks/pkg/request"
	"github.com/Taks/pkg/response"
	"github.com/go-chi/chi"
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
type TaskResponse struct {
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
		data := TaskResponse{
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

// --------------------- HANDLER DE UPDATE ---------------------
func (t *TaskHandler) UpdateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request

		// Paso 1: Leer el id de la URL y convertirlo a entero
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			//Otra forma de enviar una respuesta HTTP con error en formato texto
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// Paso 2: Leer el cuerpo de la solicitud y decodificarlo
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Paso 3: Decodificar el cuerpo y crear un map[string]any
		bodyMap := map[string]any{}
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Paso 4: Validar que todos los campos esten completos
		if err := tools.CheckFieldExistance(bodyMap, "tittle", "description", "done"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.Text(w, http.StatusBadRequest, fmt.Sprintf("%s is required", fieldError.Field))
				return
			}

			response.Text(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// Paso 5: Decodificar el JSON del cuerpo de la solicitud y asignarlo a la estructura TaskRequestBody
		var body TaskRequest
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// process
		//Paso 6: Crear una instancia de Task a partir de los datos recibidos
		task := internal.Task{
			ID:          id,
			Tittle:      body.Tittle,
			Description: body.Description,
			Done:        body.Done,
		}

		// Paso 7: Actualizar la tarea en el mapa de tareas, usando el metodo Update del repositorio
		if err := t.sv.Update(task); err != nil {
			switch {
			case errors.Is(err, internal.ErrTaskNotFound):
				response.Text(w, http.StatusNotFound, "task not found")
			case errors.Is(err, internal.ErrTaskInvalidField):
				response.Text(w, http.StatusBadRequest, "task is invalid")
			case errors.Is(err, internal.ErrTaskDuplicated):
				response.Text(w, http.StatusConflict, "task already exists")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// Paso 8: Crear  una tarea en formatoJSON que se va a enviar como respuesta del handler
		data := TaskResponse{
			ID:          task.ID,
			Tittle:      task.Tittle,
			Description: task.Description,
			Done:        task.Done,
		}

		// Paso 9: Enviar una respuesta HTTP exitosa (200 OK) junto con los datos de la tarea actualizada
		response.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "task updated",
			"data":    data,
		})
	}
}

// --------------------- HANDLER DE UPDATEPARTIAL ---------------------
func (d *TaskHandler) UpdatePartialTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// Paso 1: Leer el id de la URL y convertirlo a entero
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// Paso 2: Leer el cuerpo de la solicitud y decodificarlo
		bodyMap := make(map[string]any)
		if err := request.RequestJSON(r, &bodyMap); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// process
		// Paso 3: Actualizar la tarea en el mapa de tareas, usando el metodo UpdatePartial del repositorio
		if err := d.sv.UpdatePartial(id, bodyMap); err != nil {
			switch {
			case errors.Is(err, internal.ErrTaskNotFound):
				response.Text(w, http.StatusNotFound, "task not found")
			case errors.Is(err, internal.ErrTaskInvalidField):
				response.Text(w, http.StatusBadRequest, "task is invalid")
			case errors.Is(err, internal.ErrTaskDuplicated):
				response.Text(w, http.StatusConflict, "task already exists")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// Paso 4: Enviar una respuesta HTTP exitosa (200 OK) junto con un mensaje de tarea actualizada
		response.Text(w, http.StatusOK, "task updated")
	}
}

// --------------------- HANDLER DE DELETE ---------------------
func (d *TaskHandler) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// Paso 1: Leer el id de la URL y convertirlo a entero
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// Paso 2: Eliminar la tarea del mapa de tareas, usando el metodo Delete del repositorio
		if err := d.sv.Delete(id); err != nil {
			switch {
			case errors.Is(err, internal.ErrTaskNotFound):
				response.Text(w, http.StatusNotFound, "task not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// Paso 3: Enviar una respuesta HTTP exitosa (204 No Content) sin ningun contenido
		response.Text(w, http.StatusNoContent, "Tarea eliminada con exito")
	}
}

// --------------------- HANDLER DE GETBYID ---------------------
func (d *TaskHandler) GetTaskByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// Paso 1: Leer el id de la URL y convertirlo a entero
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		// Paso 2: Obtener la tarea del mapa de tareas, usando el metodo GetByID del repositorio
		task, err := d.sv.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrTaskNotFound):
				response.Text(w, http.StatusNotFound, "task not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// Paso 3: Crear  una tarea en formatoJSON que se va a enviar como respuesta del handler
		data := TaskResponse{
			ID:          task.ID,
			Tittle:      task.Tittle,
			Description: task.Description,
			Done:        task.Done,
		}

		// Paso 4: Enviar una respuesta HTTP exitosa (200 OK) junto con los datos de la tarea
		response.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "task found",
			"data":    data,
		})
	}
}
