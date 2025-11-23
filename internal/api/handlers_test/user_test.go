package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSetIsActive(t *testing.T) {
	r := newTestRouter()
	body := `{"user_id":"123","is_active":true}`
	req := httptest.NewRequest("POST", "/users/setIsActive", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserGetReview(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest("GET", "/users/getReview?user_id=123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
