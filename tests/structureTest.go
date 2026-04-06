package sutests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-faker/faker/v4"
)

func FillMock[T any]() (T, error) {
	var item T

	if err := faker.FakeData(&item); err != nil {
		return item, err
	}

	return item, nil

}

func StructureTest[T any]() error {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	item, err := FillMock[T]()
	if err != nil {
		return err
	}

	body, err := json.Marshal(item)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}

	response, err := httpClient.Post(server.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("criar falhou: esperado %d, obtido %d", http.StatusCreated, response.StatusCode)
	}

	response, err = httpClient.Get(server.URL)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ler falhou: esperado %d, obtido %d", http.StatusOK, response.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodPut, server.URL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	response, err = httpClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("atualizar falhou: esperado %d, obtido %d", http.StatusOK, response.StatusCode)
	}

	req, _ = http.NewRequest(http.MethodDelete, server.URL, nil)
	response, err = httpClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("deletar falhou: esperado %d, obtido %d", http.StatusNoContent, response.StatusCode)
	}

	return nil
}
