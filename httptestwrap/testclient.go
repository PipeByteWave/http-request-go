package httptestwrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/gin-gonic/gin"
)

func NewRequestTestGo(r *gin.Engine, method, url string, bod any, response any) (*httptest.ResponseRecorder, error) {
	if reflect.TypeOf(response).Kind() != reflect.Ptr {
		return nil, errors.New("error: response debe ser un puntero a una estructura o slice")
	}

	var jsonData []byte
	var err error
	if bod != nil {
		jsonData, err = json.Marshal(bod)
		if err != nil {
			return nil, fmt.Errorf("error al serializar el cuerpo de la solicitud: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error al crear la solicitud HTTP: %w", err)
	}

	if len(jsonData) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), response)
	if err != nil {
		return nil, fmt.Errorf("error al deserializar la respuesta JSON: %w", err)
	}

	return rr, nil
}
