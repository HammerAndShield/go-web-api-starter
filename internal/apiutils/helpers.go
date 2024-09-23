package apiutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-web-api-starter/internal/validator"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type Envelope map[string]any

func WriteJson(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown field %q", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger then %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func Background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error in goroutine: %v", err)
			}
		}()

		fn()
	}()
}

func BackgroundWg(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error in goroutine: %v", err)
			}
		}()

		fn()
		wg.Done()
	}()
}

func BackgroundErrGroup(eg *errgroup.Group, fn func() error) {
	eg.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error in goroutine: %v", err)
			}
		}()
		err := fn()
		if err != nil {
			return err
		}
		return nil
	})
}

func ReadStringPath(r *http.Request, key string, defaultValue string) string {
	s := r.PathValue(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func ReadIntPath(r *http.Request, key string, defaultValue int, v *validator.Validator) int {
	s := r.PathValue(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func ReadIdPath(r *http.Request, key string, v *validator.Validator) int {
	id := ReadIntPath(r, key, 0, v)
	if id <= 0 {
		v.AddError(key, "must be a positive integer value greater than 0")
	}
	return id
}

func ReadUUIDPath(r *http.Request, key string, v *validator.Validator) uuid.UUID {
	s := ReadStringPath(r, key, "")
	if s == "" {
		v.AddError(key, "must be provided")
		return uuid.Nil
	}

	id, err := uuid.Parse(s)
	if err != nil {
		v.AddError(key, "must be a valid UUID")
		return uuid.Nil
	}
	if id == uuid.Nil {
		v.AddError(key, "must be a valid UUID")
		return uuid.Nil
	}

	return id
}

func ReadCSVQuery(values url.Values, key string, defaultValue []string) ([]string, bool) {
	csv := values.Get(key)

	if csv == "" {
		return defaultValue, false
	}

	return strings.Split(csv, ","), true
}

func ReadCSVInt(values url.Values, key string, defaultValue []int, v *validator.Validator) ([]int, bool) {
	strs, exists := ReadCSVQuery(values, key, nil)
	if !exists {
		return defaultValue, false
	}

	ints := make([]int, len(strs))
	for i, str := range strs {
		id, err := strconv.Atoi(str)
		if err != nil {
			v.AddError(key, "must contain only integer values")
			return defaultValue, true
		}
		ints[i] = id
	}
	return ints, true
}

func ReadStringQuery(values url.Values, key string, defaultValue string) (string, bool) {
	s := values.Get(key)

	if s == "" {
		return defaultValue, false
	}

	return s, true
}

func ReadIntQuery(values url.Values, key string, defaultValue int, v *validator.Validator) (int, bool) {
	s := values.Get(key)

	if s == "" {
		return defaultValue, false
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue, true
	}

	return i, true
}
