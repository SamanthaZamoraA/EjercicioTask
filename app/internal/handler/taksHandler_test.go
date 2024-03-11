package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Taks/internal"
	"github.com/Taks/internal/handler"
	"github.com/Taks/internal/repository"
	"github.com/Taks/internal/service"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

// Test de handler GetByID
func TestGetByID(t *testing.T) {

	//Test obtener de forma correcta una tarea
	t.Run("Success - GetById Tak", func(t *testing.T) {

		//arrange

		//Hacer una tarea de prueba
		db := map[int]internal.Task{
			1: {
				ID:          1,
				Tittle:      "task 1",
				Description: "description 1",
				Done:        false,
			},
		}

		//Inicializar dependencias
		//Dependencia del repository
		rp := repository.NewTaskMap(db, 0)

		//Dependencia del service
		sv := service.NewTaskService(rp)

		//Dependencia del handler
		h := handler.NewTaskHandler(sv)

		//Declarar el método a testear
		hdFunc := h.GetTaskByID()

		//act

		//Hacer el request
		req := httptest.NewRequest("GET", "/task/getById/1", nil)
		//Hacer el contexto
		chiCtx := chi.NewRouteContext()
		//Pasar el id
		chiCtx.URLParams.Add("id", "1")
		//Replantear el request
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		//Hacer el responde
		res := httptest.NewRecorder()

		//Pasarlo al handler
		hdFunc(res, req)

		//assert

		//Lo que se espera
		expectedTask := map[string]any{
			"message": "task found",
			"data": handler.TaskResponse{
				ID:          1,
				Tittle:      "task 1",
				Description: "description 1",
				Done:        false,
			},
		}

		//Parsear el JSON y validar el error
		expectedTaskJSON, err := json.Marshal(expectedTask)
		require.NoError(t, err)

		//Hacer las comparaciones
		require.Equal(t, http.StatusOK, res.Code)
		require.JSONEq(t, string(expectedTaskJSON), res.Body.String())
		require.True(t, strings.HasPrefix(res.Header().Get("Content-Type"), "application/json"))
	})
}
