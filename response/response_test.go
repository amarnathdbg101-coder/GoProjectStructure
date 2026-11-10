package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amarnathdbg101-coder/GoProjectStructure/response"
)

func TestResponse_Success(t *testing.T) {
	w := httptest.NewRecorder()
	response.Success(w, "operation completed", map[string]string{"foo": "bar"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success || res.Message != "operation completed" {
		t.Fatalf("unexpected response payload: %+v", res)
	}
}

func TestResponse_Error(t *testing.T) {
	w := httptest.NewRecorder()
	response.Error(w, http.StatusBadRequest, "invalid payload")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var res response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Success || res.Error != "invalid payload" {
		t.Fatalf("unexpected error response payload: %+v", res)
	}
}
