package sutests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
)

func FillMock[T any]() (T, error) {
	var item T

	if err := faker.FakeData(&item); err != nil {
		return item, err
	}

	return item, nil

}

type envelopeResponse struct {
	Message string          `json:"message"`
	Error   string          `json:"error"`
	Data    json.RawMessage `json:"data"`
}

func StructureTest[T any](app *fiber.App, baseRoute string) error {
	if app == nil {
		return fmt.Errorf("fiber app cannot be nil")
	}

	baseRoute = strings.TrimSpace(baseRoute)
	if baseRoute == "" {
		return fmt.Errorf("base route cannot be empty")
	}

	if !strings.HasPrefix(baseRoute, "/") {
		baseRoute = "/" + baseRoute
	}

	item, err := FillMock[T]()
	if err != nil {
		return err
	}

	body, err := json.Marshal(item)
	if err != nil {
		return err
	}

	response, err := app.Test(newJSONRequest(http.MethodPost, baseRoute, body))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("criar falhou: esperado %d, obtido %d", http.StatusCreated, response.StatusCode)
	}

	created, err := parseData[T](response)
	if err != nil {
		return fmt.Errorf("falha ao parsear resposta de create: %w", err)
	}

	id, err := extractUintID(created)
	if err != nil {
		return fmt.Errorf("falha ao extrair ID na resposta de create: %w", err)
	}

	response, err = app.Test(newRequest(http.MethodGet, fmt.Sprintf("%s/%d", baseRoute, id)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ler por id falhou: esperado %d, obtido %d", http.StatusOK, response.StatusCode)
	}

	readByID, err := parseData[T](response)
	if err != nil {
		return fmt.Errorf("falha ao parsear resposta de get by id: %w", err)
	}

	readByIDValue := reflect.ValueOf(readByID)
	createdValue := reflect.ValueOf(created)
	if readByIDValue != createdValue {
		if _, idErr := extractUintID(readByID); idErr != nil {
			return fmt.Errorf("payload de get by id invalido: %w", idErr)
		}
	}

	response, err = app.Test(newRequest(http.MethodGet, baseRoute))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("listar falhou: esperado %d, obtido %d", http.StatusOK, response.StatusCode)
	}

	list, err := parseData[[]T](response)
	if err != nil {
		return fmt.Errorf("falha ao parsear resposta de get all: %w", err)
	}

	if len(list) == 0 {
		return fmt.Errorf("listar falhou: nenhum item retornado")
	}

	found := false
	for _, current := range list {
		currentID, idErr := extractUintID(current)
		if idErr == nil && currentID == id {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("listar falhou: item criado nao encontrado")
	}

	updated := created
	setFirstStringField(&updated, "updated")

	setIDErr := setUintID(&updated, id)
	if setIDErr != nil {
		return fmt.Errorf("falha ao preparar payload de update: %w", setIDErr)
	}

	updatedBody, err := json.Marshal(updated)
	if err != nil {
		return err
	}

	response, err = app.Test(newJSONRequest(http.MethodPut, baseRoute, updatedBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("atualizar falhou: esperado %d, obtido %d", http.StatusOK, response.StatusCode)
	}

	_, err = parseData[T](response)
	if err != nil {
		return fmt.Errorf("falha ao parsear resposta de update: %w", err)
	}

	response, err = app.Test(newRequest(http.MethodDelete, fmt.Sprintf("%s/%d", baseRoute, id)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("deletar falhou: esperado %d, obtido %d", http.StatusOK, response.StatusCode)
	}

	response, err = app.Test(newRequest(http.MethodGet, fmt.Sprintf("%s/%d", baseRoute, id)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("pos-delete get by id falhou: esperado %d, obtido %d", http.StatusNotFound, response.StatusCode)
	}

	return nil
}

func newRequest(method, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	return req
}

func newJSONRequest(method, url string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func parseData[T any](response *http.Response) (T, error) {
	var output T

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return output, err
	}

	var envelope envelopeResponse
	if err := json.Unmarshal(body, &envelope); err != nil {
		return output, err
	}

	if len(envelope.Data) == 0 {
		return output, fmt.Errorf("resposta sem campo data")
	}

	if err := json.Unmarshal(envelope.Data, &output); err != nil {
		return output, err
	}

	return output, nil
}

func extractUintID[T any](item T) (uint, error) {
	value := reflect.ValueOf(item)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return 0, fmt.Errorf("item nil")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return 0, fmt.Errorf("item nao eh struct")
	}

	idField := value.FieldByName("ID")
	if !idField.IsValid() {
		idField = value.FieldByName("Id")
	}

	if !idField.IsValid() {
		return 0, fmt.Errorf("campo ID/Id nao encontrado")
	}

	switch idField.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		id := uint(idField.Uint())
		if id == 0 {
			return 0, fmt.Errorf("id retornou zero")
		}
		return id, nil
	default:
		return 0, fmt.Errorf("campo ID/Id nao eh inteiro sem sinal")
	}
}

func setUintID[T any](item *T, id uint) error {
	if item == nil {
		return fmt.Errorf("item nil")
	}

	value := reflect.ValueOf(item)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return fmt.Errorf("item deve ser ponteiro")
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("item deve ser struct")
	}

	idField := value.FieldByName("ID")
	if !idField.IsValid() {
		idField = value.FieldByName("Id")
	}

	if !idField.IsValid() || !idField.CanSet() {
		return fmt.Errorf("campo ID/Id nao encontrado ou nao setavel")
	}

	switch idField.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		idField.SetUint(uint64(id))
		return nil
	default:
		return fmt.Errorf("campo ID/Id nao eh inteiro sem sinal")
	}
}

func setFirstStringField[T any](item *T, suffix string) {
	if item == nil {
		return
	}

	value := reflect.ValueOf(item)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			current := field.String()
			if strings.TrimSpace(current) == "" {
				field.SetString(suffix)
			} else {
				field.SetString(current + "-" + suffix)
			}
			return
		}
	}
}
