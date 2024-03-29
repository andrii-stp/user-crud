package handler_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrii-stp/users-crud/router"
	"github.com/andrii-stp/users-crud/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Handler Suite")
}

func ExecuteRequest(logger *slog.Logger, req *http.Request, repo storage.UserRepository) *httptest.ResponseRecorder {
	r := router.Router(logger, repo)
	nr := httptest.NewRecorder()

	r.ServeHTTP(nr, req)

	return nr
}

func Deserialize(d string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(d), &m); err != nil {
		return nil, fmt.Errorf("failed to deserialize. %w", err)
	}

	return m, nil
}

func DeserializeList(d string) ([]map[string]interface{}, error) {
	var l []map[string]interface{}
	if err := json.Unmarshal([]byte(d), &l); err != nil {
		return nil, fmt.Errorf("failed to deserialize a list. %w", err)
	}

	return l, nil
}
