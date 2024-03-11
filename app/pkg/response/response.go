package response

import (
	"encoding/json"
	"net/http"
)

/*
Esta funcion es para escribir respuestas en formato texto:

- > w http.ResponseWriter: Un objeto que permite escribir una respuesta HTTP al cliente.
- > code int: El código de estado HTTP que se enviará en la respuesta.
- > body string: El cuerpo del mensaje que se enviará como respuesta.
*/
func Text(w http.ResponseWriter, code int, body string) {

	//Configura el encabezado Content-Type de la respuesta HTTP para indicar que la respuesta será texto plano codificado en UTF-8.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	//Establece el código de estado HTTP de la respuesta en el valor proporcionado por code
	w.WriteHeader(code)

	/*
		Toma el contenido de la variable body y lo envía como respuesta al cliente que hizo la solicitud HTTP:

		- >[]byte(body): Convierte el contenido de body (que es una cadena de texto) en un formato de bytes que se puede enviar a través de la red.s
		- > w.Write([]byte(body)): Envía esos bytes como parte del cuerpo de la respuesta HTTP al cliente que hizo la solicitud.
	*/
	w.Write([]byte(body))
}

/*
Esta funcion es para escribir respuestas en formato JSON:

  - > w http.ResponseWriter: Un objeto que permite escribir una respuesta HTTP al cliente.
  - > code int: El código de estado HTTP que se enviará en la respuesta.
  - > body any: El cuerpo del mensaje JSON que se enviará como respuesta. El tipo any puede ser de cualquier tipo.
*/
func ResponseJSON(w http.ResponseWriter, code int, body any) {

	//Configura el encabezado Content-Type de la respuesta HTTP para indicar que la respuesta será JSON codificado en UTF-8.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//Establece el código de estado HTTP de la respuesta en el valor proporcionado por code
	w.WriteHeader(code)

	//Escribe el cuerpo de la respuesta JSON en la respuesta HTTP
	if err := json.NewEncoder(w).Encode(body); err != nil {
		//Si ocurre un error codificando el JSON, enviar una respuesta HTTP con error
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
}
