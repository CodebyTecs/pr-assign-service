package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPRCreate(t *testing.T) {
	r := newTestRouter()
	body := `{
		"pull_request_id": "1",
		"pull_request_name": "test-pr",
		"author_id": "1" 
	}`
	req := httptest.NewRequest("POST", "/pullRequest/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPRReassign(t *testing.T) {
	r := newTestRouter()
	body := `{
		"pull_request_id": "1",
		"old_user_id": "1"
	}`
	req := httptest.NewRequest("POST", "/pullRequest/reassign", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPRMerge(t *testing.T) {
	r := newTestRouter()
	body := `{"pull_request_id": "1"}`
	req := httptest.NewRequest("POST", "/pullRequest/merge", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
