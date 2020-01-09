package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// ParseJSONRequest tries to parse a JSON request body into target struct. In case of error, it responds with the proper error and returns false. Based on https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func ParseJSONRequest(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(target)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var msg string

		switch {
		case errors.As(err, &syntaxError):
			msg = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg = fmt.Sprintf("Request body contains badly-formed JSON")
			w.WriteHeader(http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			msg = fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			w.WriteHeader(http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg = fmt.Sprintf("Request body contains unknown field %s", fieldName)
			w.WriteHeader(http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			msg = "Request body must not be empty"
			w.WriteHeader(http.StatusBadRequest)

		default:
			log.Println(err.Error())
			msg = http.StatusText(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": msg})

		return false
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": msg})
		return false
	}

	return true
}
