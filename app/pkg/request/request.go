package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Variable global para capturar errores
var (
	ErrRequestJSONInvalid      = errors.New("Solicitud JSON invalida")
	ErrRequestPathParamInvalid = errors.New("Parametro de ruta invalido")
)

/*
Recibe:
  - > r *http.Request: puntero a una estructura de solicitud HTTP (http.Request).
    Es la solicitud que se está procesando y de la cual se espera extraer datos en formato JSON.

- > ptr any: puntero a cualquier tipo de dato. Es el puntero al cual se le asignará el valor extraído del JSON.

Devuelve:
- > err error: error. Es el error que se produjo durante la extracción del JSON.
*/
func RequestJSON(r *http.Request, ptr any) (err error) {

	/*
		-> NewDecoder(r.Body): crea un nuevo decodificador JSON a partir del cuerpo de la solicitud r.Body.
		-> Decode(ptr): decodifica el JSON del cuerpo de la solicitud r.Body y lo asigna al puntero ptr.
	*/
	err = json.NewDecoder(r.Body).Decode(ptr)

	/* Si hay un error durante la decodificación del JSON, crea un nuevo error que incluye el error original
	y un mensaje de error personalizado.
	*/
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrRequestJSONInvalid, err)
		return
	}
	return
}

/*
Función:
- > PathLastParam: extrae el último parámetro de la ruta de la solicitud HTTP.

Recibe:
  - > r *http.Request: puntero a una estructura de solicitud HTTP (http.Request).
    Es la solicitud que se está procesando y de la cual se espera extraer datos en formato JSON.

Devuelve:
- > value string: Es el último parámetro de la ruta.
- > err error: error. Es el error que se produjo durante la extracción del JSON.
*/
func PathLastParam(r *http.Request) (value string, err error) {

	// Contiene la ruta de la solicitud HTTP
	path := r.URL.Path

	// regexp.MustCompile: valida la forma de la ruta y devuelve un error si no es válida.
	rx := regexp.MustCompile(`^/(.*/)*([0-9a-zA-Z]+)$`)

	// MatchString: compara la ruta (path) con la expresión regular (rx) y devuelve true si es válida.
	if !rx.MatchString(path) {
		err = ErrRequestPathParamInvalid
		return
	}

	/*
		-> strings.Split(path, "/"): Divide la cadena path en segmentos utilizando el carácter "/" como delimitador.
		Cada segmento se almacenará en un slice de strings.

		-> sl[len(sl)-1]: Accede al último elemento del slice sl. Esto se hace usando len(sl)-1 como índice para
		acceder al último elemento del slice.

		value = sl[len(sl)-1]: Asigna el último elemento del slice (que es el último segmento de la URL) a la variable value.

		return: Devuelve value como el último parámetro de la URL.
	*/
	sl := strings.Split(path, "/")
	value = sl[len(sl)-1]
	return
}
