package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Vesino/todos/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	if status != 200 && status != 201 {
		w.WriteHeader(status)
	}

	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body into the target destination.
	err := dec.Decode(dst)

	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is()
		// and return a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error
		// to our handler. At the end of this chapter we'll talk about panicking
		// versus returning errors, and discuss why it's an appropriate thing to do in
		// this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) background(fn func()) {

	app.wg.Add(1)

	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
	}()
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}

	return i
}
