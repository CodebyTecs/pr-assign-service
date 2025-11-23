package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamAdd(t *testing.T) {
	r := newTestRouter()
	body := `{
		"team_name": "test",
		"members": [
			{"user_id": "1", "username": "test1", "is_active": true},
			{"user_id": "2", "username": "test2", "is_active": true}
		]
}`
	req := httptest.NewRequest("POST", "/team/add", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTeamGet(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest("GET", "/team/get?team_name=backend", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
